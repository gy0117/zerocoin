package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/withdraw"
	"time"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/model"
	"ucenter-rpc/internal/svc"
	"zero-common/kafka"
	"zero-common/operate"
	"zero-common/tools"
	"zero-common/zerodb"
	"zero-common/zerodb/tran"
	"zero-common/zerr"
)

var ErrWithdraw = zerr.NewCodeErr(zerr.WITHDRAW_ERROR)
var ErrFindWithdrawRecord = zerr.NewCodeErr(zerr.WITHDRAW_FIND_RECORD)

const withdrawVerifyCode = "WITHDRAW::VERIFY::"
const topicBtcWithdraw = "btc_withdraw"

type WithdrawLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	mAddressDomain     *domain.UserAddressDomain
	userDomain         *domain.UserDomain
	memberWalletDomain *domain.WalletDomain
	withdrawDomain     *domain.WithdrawDomain
	transaction        tran.Transaction
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		mAddressDomain:     domain.NewUserAddressDomain(svcCtx.DB),
		userDomain:         domain.NewUserDomain(svcCtx.DB),
		memberWalletDomain: domain.NewWalletDomain(svcCtx.DB),
		withdrawDomain:     domain.NewWithdrawDomain(svcCtx.DB),
		transaction:        tran.NewTransaction(svcCtx.DB.Conn),
	}
}

func (l *WithdrawLogic) FindAddressesByCoinId(in *withdraw.WithdrawRequest) (*withdraw.AddressSimpleList, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	addresses, err := l.mAddressDomain.FindAddressesByCoinId(ctx, in.UserId, in.CoinId)
	if err != nil {
		return nil, errors.Wrapf(ErrGetAddress, "获取地址失败, uid: %d, coin: %d, err: %v", in.UserId, in.CoinId, err)
	}

	list := make([]*withdraw.AddressSimple, len(addresses))
	for i, v := range addresses {
		item := &withdraw.AddressSimple{}
		item.Address = v.Address
		item.Remark = v.Remark

		list[i] = item
	}

	resp := &withdraw.AddressSimpleList{
		List: list,
	}
	return resp, nil
}

func (l *WithdrawLogic) SendCode(in *withdraw.SendCodeReq) (*withdraw.EmptyResp, error) {
	// 这里就随机生成
	code := tools.GenerateVerifyCode()
	logx.Infof("withdraw verifyCode: %s\n", code)

	go func() {
		// 这里应该发送验证码
	}()

	key := withdrawVerifyCode + in.Phone
	err := l.svcCtx.Cache.SetWithExpireCtx(context.Background(), key, code, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	return &withdraw.EmptyResp{}, nil
}

// Withdraw
// 1. 短信验证码校验
// 2. 校验交易密码是否正确
// 3. 根据用户id和unit，查询用户钱包，判断余额是否足够
// 4. 冻结用户的钱 提现币 经过比特币网络需要时间，因此需要冻结
// 5. 记录用户的提现
// 6. 发送用户的提现事件到MQ，MQ消费者去处理提现（创建交易 广播到比特币网络）
func (l *WithdrawLogic) Withdraw(in *withdraw.WithdrawRequest) (*withdraw.EmptyResp, error) {
	member, err := l.userDomain.FindUserById(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(ErrFindUser, "查找用户失败, uid: %d, err: %v", in.UserId, err)
	}
	var codeStr string
	key := withdrawVerifyCode + member.MobilePhone
	err = l.svcCtx.Cache.Get(key, &codeStr)
	if err != nil {
		return nil, err
	}
	if in.Code != codeStr {
		// 验证码输入有误
		return nil, errors.Wrapf(ErrWithdraw, "验证码输入错误, uid: %d, err: %v", in.UserId, err)
	}

	if in.JyPassword != member.JyPassword {
		// 交易密码错误
		return nil, errors.Wrapf(ErrWithdraw, "交易密码输入错误, uid: %d", in.UserId)
	}

	memberWallet, err := l.memberWalletDomain.FindWalletByMemIdAndCoinName(l.ctx, in.UserId, in.Unit)
	if err != nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找钱包失败, uid: %d, err: %v", in.UserId, err)
	}
	if memberWallet.Balance < in.Amount {
		// 如果钱包中的余额（这里的余额就是coinName对应的数量）
		return nil, errors.Wrapf(ErrWithdraw, "余额不足, uid: %d, balance: %f", in.UserId, memberWallet.Balance)
	}

	// 下面是一个事务
	err = l.transaction.Action(func(conn zerodb.DbConn) error {
		// 冻结
		if err := l.memberWalletDomain.FreezeWithConn(l.ctx, conn, in.UserId, in.Amount, in.Unit); err != nil {
			return errors.Wrapf(ErrWithdraw, "冻结钱包失败, uid: %d, amount: %f, unit: %s, err: %v", in.UserId, in.Amount, in.Unit, err)
		}

		// 保存记录
		// 保存提现记录
		record := &model.WithdrawRecord{}
		record.CoinId = memberWallet.CoinId
		record.Address = in.Address
		record.Fee = in.Fee
		record.TotalAmount = in.Amount
		record.ArrivedAmount = operate.SubFloor(in.Amount, in.Fee, 10)
		record.Remark = ""
		record.CanAutoWithdraw = 0
		record.IsAuto = 0
		record.Status = 0 //审核中
		record.CreateTime = time.Now().UnixMilli()
		record.DealTime = 0
		record.UserId = in.UserId
		record.TransactionNumber = "" //目前还没有交易编号

		if err := l.withdrawDomain.SaveRecord(l.ctx, conn, record); err != nil {
			return errors.Wrapf(ErrWithdraw, "保存提现记录失败, uid: %d, err: %v", in.UserId, err)
		}

		// 发送MQ
		marshal, _ := json.Marshal(record)

		data := kafka.KafkaData{
			Topic: topicBtcWithdraw,
			Key:   []byte(fmt.Sprintf("%d", in.UserId)),
			Data:  marshal,
		}
		for i := 0; i < 3; i++ {
			err = l.svcCtx.KCli.SendSync(data)
			if err != nil {
				time.Sleep(time.Millisecond * 250)
				continue
			}
			break
		}
		if err != nil {
			return errors.Wrapf(ErrWithdraw, "发送kafka失败, uid: %d, err: %v", in.UserId, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &withdraw.EmptyResp{}, nil
}

func (l *WithdrawLogic) WithdrawRecord(req *withdraw.WithdrawRequest) (*withdraw.RecordList, error) {
	list, total, err := l.withdrawDomain.WithdrawRecord(l.ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.Wrapf(ErrFindWithdrawRecord, "查询提现记录失败, uid: %d, err: %v", req.UserId, err)
	}

	voList := make([]*model.WithdrawRecordVo, len(list))
	for i, v := range list {
		coin, err := l.svcCtx.MarketRpc.FindCoinByCoinId(l.ctx, &market.MarketRequest{
			CoinId: v.CoinId,
		})
		if err != nil {
			continue
		}
		voList[i] = v.ToVo(coin)
	}

	var res []*withdraw.WithdrawRecord
	_ = copier.Copy(&res, voList)
	return &withdraw.RecordList{
		List:  res,
		Total: total,
	}, nil
}
