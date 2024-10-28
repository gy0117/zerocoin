package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/wallet"
	"time"
	"ucenter-rpc/internal/config"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/svc"
	"ucenter-rpc/internal/verify"
	"zero-common/btc"
	"zero-common/dtmutil"
	"zero-common/operate"
	"zero-common/tools"
	"zero-common/zerr"
)

var ErrFindWallet = zerr.NewCodeErr(zerr.FIND_WALLET_ERROR)
var ErrResetWalletAddress = zerr.NewCodeErr(zerr.RESET_WALLET_ADDRESS_ERROR)
var ErrGetTransactions = zerr.NewCodeErr(zerr.GET_TRANSACTIONS_ERROR)
var ErrGetAddress = zerr.NewCodeErr(zerr.GET_ADDRESS_ERROR)

type WalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	captchaVerify *verify.MachineVerify
	userDomain    *domain.UserDomain
	walletDomain  *domain.WalletDomain
	mtDomain      *domain.TransactionDomain
}

func NewWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletLogic {
	return &WalletLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		captchaVerify: verify.NewMachineVerify(),
		userDomain:    domain.NewUserDomain(svcCtx.DB),
		walletDomain:  domain.NewWalletDomain(svcCtx.DB),
		mtDomain:      domain.NewTransactionDomain(svcCtx.DB),
	}
}

func (l *WalletLogic) FindWalletBySymbol(in *wallet.WalletReq) (*wallet.UserWallet, error) {
	ctx := context.Background()
	// 1. 根据coinName，在表coin中，找到coin信息（market rpc服务中就有）
	coinResp, err := l.svcCtx.MarketRpc.FindCoinInfo(ctx, &market.MarketRequest{
		Unit: in.CoinName,
	})

	if err != nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找coin表失败，coin: %s, err: %v", in.CoinName, err)
	}
	if coinResp == nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找coin表失败，数据为空，coin: %s", in.CoinName)
	}

	// 2. 根据coin_name，在表user_wallet中，找到row
	memberWalletCoin, err := l.walletDomain.FindWallet(ctx, in.UserId, in.CoinName, coinResp)
	if err != nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找user_wallet表失败，coin: %s, err: %v", in.CoinName, err)
	}
	resp := &wallet.UserWallet{}
	_ = copier.Copy(resp, memberWalletCoin)
	return resp, nil
}

func (l *WalletLogic) FindWallet(in *wallet.WalletReq) (*wallet.FindWalletResp, error) {
	// 从user_wallet中找到对应userId的钱包信息
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	memberWallets, err := l.walletDomain.FindWalletsByUserId(ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找user_wallet表失败，uid: %s, err: %v", in.UserId, err)
	}

	// 从redis中获取汇率
	var cnyRateStr string
	var cnyRate float64 = 7
	_ = l.svcCtx.Cache.Get("USDT::CNY::RATE", &cnyRateStr)
	if cnyRateStr != "" {
		cnyRate = tools.ToFloat64(cnyRateStr)
	}

	list := make([]*wallet.UserWallet, len(memberWallets))
	index := 0
	for _, v := range memberWallets {
		mwc := &wallet.UserWallet{}
		_ = copier.Copy(mwc, v)

		coinResp, err := l.findCoinByUnit(ctx, v.CoinName)
		if err != nil {
			logx.Error(err)
			list[index] = mwc
			index++
			continue
		}
		walletCoin := &wallet.Coin{}
		_ = copier.Copy(walletCoin, coinResp)
		mwc.Coin = walletCoin

		if v.CoinName == "USDT" {
			mwc.Coin.UsdRate = 1
			mwc.Coin.CnyRate = cnyRate
		} else {
			var rateStr string
			var rate float64 = 10000
			_ = l.svcCtx.Cache.Get(v.CoinName+"::USDT::RATE", &rateStr)
			if rateStr != "" {
				rate = tools.ToFloat64(rateStr)
			}
			// 1 BTC 对应 rate个 USDT，那么 1 BTC 对应 cnyRate * rate个 USDT
			mwc.Coin.UsdRate = rate
			mwc.Coin.CnyRate = operate.MulFloor(cnyRate, rate, 8)
		}

		logx.Info("%s | UsdRate: %f, CnyRate: %f\n", v.CoinName, mwc.Coin.UsdRate, mwc.Coin.CnyRate)

		list[index] = mwc
		index++
	}
	return &wallet.FindWalletResp{
		List: list,
	}, nil
}

func (l *WalletLogic) findCoinByUnit(ctx context.Context, coinName string) (*mclient.CoinResp, error) {
	marketReq := &market.MarketRequest{
		Unit: coinName,
	}
	coinResp, err := l.svcCtx.MarketRpc.FindCoinInfo(ctx, marketReq)
	if err != nil {
		return nil, err
	}
	if coinResp == nil {
		return nil, err
	}
	return coinResp, nil
}

func (l *WalletLogic) ResetWalletAddress(in *wallet.WalletReq) (*wallet.WalletResp, error) {
	userWallet, err := l.walletDomain.FindWalletByMemIdAndCoinName(l.ctx, in.UserId, in.CoinName)
	if err != nil {
		return nil, errors.Wrapf(ErrFindWallet, "查找用户失败，uid: %s, coin: %s, err: %v", in.UserId, in.CoinName, err)
	}
	if in.CoinName == "BTC" {
		if userWallet.Address == "" {
			// 生成地址
			newWallet, err := btc.NewWallet()
			if err != nil {
				return nil, errors.Wrapf(ErrResetWalletAddress, "重置钱包地址失败，uid: %d, coin: %s, err: %v", in.UserId, in.CoinName, err)
			}
			address := newWallet.GenerateBitcoinTestAddress()
			privateKey := newWallet.GenerateBitcoinPrivateKey()

			userWallet.Address = string(address)
			userWallet.AddressPrivateKey = privateKey[:50]

			if err := l.walletDomain.UpdateWalletAddress(l.ctx, userWallet); err != nil {
				return nil, errors.Wrapf(ErrResetWalletAddress, "更新钱包地址失败 uid: %d, coin: %s, err: %v", in.UserId, in.CoinName, err)
			}
		}
	}
	return &wallet.WalletResp{}, nil
}

func (l *WalletLogic) GetAllTransactions(in *wallet.AssetReq) (*wallet.UserTransactionListResp, error) {
	transactionVos, total, err := l.mtDomain.GetTransactions(l.ctx, in.PageNo, in.PageSize, in.UserId, in.Symbol, in.StartTime, in.EndTime, in.Type)
	if err != nil {
		return nil, errors.Wrapf(ErrGetTransactions, "获取交易记录失败, uid: %d, err: %v", in.UserId, err)
	}

	// rpc对象
	var list []*wallet.UserTransaction
	_ = copier.Copy(&list, transactionVos)
	resp := &wallet.UserTransactionListResp{
		List:  list,
		Total: total,
	}
	return resp, nil
}

func (l *WalletLogic) GetAddress(in *wallet.AssetReq) (*wallet.AddressListResp, error) {
	address, err := l.walletDomain.GetAddress(l.ctx, in.CoinName)
	if err != nil {
		return nil, errors.Wrapf(ErrGetAddress, "获取地址失败, uid: %d, err: %v", in.UserId, err)
	}
	return &wallet.AddressListResp{
		List: address,
	}, nil
}

func dbGet(c config.Config) *dtmutil.DB {
	var dbConf = dtmcli.DBConf{
		Driver:   c.Mysql.Driver,
		Host:     c.Mysql.Host,
		Port:     c.Mysql.Port,
		User:     c.Mysql.User,
		Password: c.Mysql.Password,
		Db:       c.Mysql.Db,
	}
	return dtmutil.DbGet(dbConf)
}

func (l *WalletLogic) FreezeUserAsset(in *wallet.FreezeUserAssetReq) (*wallet.Empty, error) {
	logx.Info("saga -> 冻结资产")
	// TODO 需要吗？1. 根据orderId，查询订单状态，如果订单交易中，说明已经被处理了  感觉是不需要了

	// 冻结
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, errors.Wrap(status.Error(codes.Aborted, err.Error()), "freeze_asset_logic create barrier failed")
	}
	tx := dbGet(l.svcCtx.Config).DB.Begin()
	//sourceTx := tx.Statement.ConnPool.(*gorm.PreparedStmtTX).Tx.(*sql.Tx)
	sourceTx := tx.Statement.ConnPool.(*sql.Tx)

	assetDomain := domain.NewAssetDomain(tx)
	err = barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
		// 这里不需要判断symbol
		err := assetDomain.Freeze(l.ctx, in.GetUid(), in.GetMoney(), in.GetSymbol())
		if err != nil {
			return fmt.Errorf("freeze user asset, req: %+v, err: %v", in, err)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &wallet.Empty{}, nil
}

// 将冻结的钱还原回去
func (l *WalletLogic) UnFreezeUserAsset(in *wallet.FreezeUserAssetReq) (*wallet.Empty, error) {
	logx.Info("saga -> 解冻资产")
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, errors.Wrap(status.Error(codes.Aborted, err.Error()), "unfreeze asset create barrier failed")
	}
	tx := dbGet(l.svcCtx.Config).DB.Begin()
	sourceTx := tx.Statement.ConnPool.(*sql.Tx)
	asssetDomain := domain.NewAssetDomain(tx)
	err = barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
		err := asssetDomain.UnFreeze(l.ctx, in.GetUid(), in.GetMoney(), in.GetSymbol())
		if err != nil {
			return fmt.Errorf("unfreeze user asset failed, req: %+v, err: %v", in, err)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &wallet.Empty{}, nil
}

// TODO
func (l *WalletLogic) DeductUserAsset(in *wallet.DeductUserAssetReq) (*wallet.Empty, error) {
	logx.Info("saga -> 扣减资产")
	return &wallet.Empty{}, nil
}

// TODO
func (l *WalletLogic) AddUserAsset(in *wallet.AddUserAssetReq) (*wallet.Empty, error) {
	logx.Info("saga -> 增加资产")
	return &wallet.Empty{}, nil
}
