package withdraw

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/user"
	"grpc-common/ucenter/types/wallet"
	"grpc-common/ucenter/types/withdraw"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/pages"
)

type WithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawLogic) GetSupportedCoinInfo() ([]*types.WithdrawWalletInfo, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)

	// 1. 查询所有coin信息（可能没必要）
	//coins, err := l.svcCtx.MarketRpc.FindAllCoins(ctx, &market.MarketRequest{})
	//if err != nil {
	//	logx.Error(err)
	//	return nil, err
	//}
	//coinMap := make(map[string]*market.Coin)
	//for _, v := range coins.List {
	//	coinMap[v.Unit] = v
	//}

	// 2. 根据用户id 查询用户的钱包信息
	memberWallets, err := l.svcCtx.UCWalletRpc.FindWallet(ctx, &wallet.WalletReq{
		UserId: userId,
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	// 3. 组装数据[]*types.WithdrawWalletInfo
	list := make([]*types.WithdrawWalletInfo, len(memberWallets.List))
	for i, v := range memberWallets.List {
		item := &types.WithdrawWalletInfo{}

		//coin := coinMap[wallet.Coin.Unit]
		coin := v.Coin
		item.Balance = v.Balance
		item.WithdrawScale = int(coin.WithdrawScale)
		item.MaxTxFee = coin.MaxTxFee
		item.MinTxFee = coin.MinTxFee
		item.MaxAmount = coin.MaxWithdrawAmount
		item.MinAmount = coin.MinWithdrawAmount
		item.Name = coin.GetName()
		item.NameCn = coin.NameCn
		item.Threshold = coin.WithdrawThreshold
		item.Unit = coin.Unit
		item.AccountType = int(coin.AccountType)
		if coin.CanAutoWithdraw == 0 {
			item.CanAutoWithdraw = "true"
		} else {
			item.CanAutoWithdraw = "false"
		}

		// 根据coin_id和user_id，找到对应的user_address
		simpleList, err := l.svcCtx.UCWithdrawRpc.FindAddressesByCoinId(ctx, &withdraw.WithdrawRequest{
			CoinId: int64(coin.Id),
			UserId: userId,
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		var as []types.AddressSimple
		err = copier.Copy(&as, simpleList.List)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		item.Addresses = as

		list[i] = item
	}

	return list, nil
}

func (l *WithdrawLogic) SendCode() error {
	// 1. 根据userId，找到手机号码
	// 2. 发短信

	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)

	memberResp, err := l.svcCtx.UserRpc.FindUserById(ctx, &user.UserRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Error(err)
		return err
	}
	phone := memberResp.MobilePhone
	_, err = l.svcCtx.UCWithdrawRpc.SendCode(ctx, &withdraw.SendCodeReq{
		Phone: phone,
	})
	if err != nil {
		logx.Error(err)
		return err
	}
	return nil
}

func (l *WithdrawLogic) Withdraw(in *types.WithdrawReq) error {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()
	userId := ctx.Value("userId").(int64)

	req := &withdraw.WithdrawRequest{
		UserId:     userId,
		Unit:       in.Unit,
		JyPassword: in.JyPassword,
		Code:       in.Code,
		Address:    in.Address,
		Amount:     in.Amount,
		Fee:        in.Fee,
	}
	_, err := l.svcCtx.UCWithdrawRpc.Withdraw(ctx, req)
	if err != nil {
		logx.Error(err)
		return err
	}
	return nil
}

func (l *WithdrawLogic) Record(req *types.WithdrawReq) (*pages.PageResult, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)
	record, err := l.svcCtx.UCWithdrawRpc.WithdrawRecord(ctx, &withdraw.WithdrawRequest{
		UserId:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]any, len(record.List))
	for i, v := range record.List {
		list[i] = v
	}
	return pages.New(list, req.Page, req.PageSize, record.Total), nil
}
