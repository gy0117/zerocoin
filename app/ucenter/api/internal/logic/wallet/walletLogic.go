package wallet

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/wallet"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/pages"
)

type WalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletLogic {
	return &WalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WalletLogic) GetWalletInfo(req *types.WalletReq) (*types.UserWallet, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// 这里的参数需要转换，api层和rpc层的对象不要用同一个
	// rpc 调用
	userId := ctx.Value("userId").(int64)
	resp, err := l.svcCtx.UCWalletRpc.FindWalletBySymbol(ctx, &wallet.WalletReq{
		CoinName: req.CoinName,
		UserId:   userId,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "uid: %d, coinName: %s", userId, req.CoinName)
	}
	result := &types.UserWallet{}
	_ = copier.Copy(result, resp)
	return result, nil
}

func (l *WalletLogic) FindWallet() ([]*types.UserWallet, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// rpc调用
	userId := ctx.Value("userId").(int64)
	in := &wallet.WalletReq{
		UserId: userId,
	}
	findWalletResp, err := l.svcCtx.UCWalletRpc.FindWallet(ctx, in)
	if err != nil {
		return nil, errors.Wrapf(err, "uid: %d", userId)
	}
	var list []*types.UserWallet

	_ = copier.Copy(&list, findWalletResp.List)
	return list, nil
}

func (l *WalletLogic) ResetWalletAddress(req *types.WalletReq) (string, error) {

	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// rpc调用
	userId := ctx.Value("userId").(int64)
	in := &wallet.WalletReq{
		UserId:   userId,
		CoinName: req.Unit,
	}
	_, err := l.svcCtx.UCWalletRpc.ResetWalletAddress(ctx, in)
	if err != nil {
		return "", errors.Wrapf(err, "uid: %d", userId)
	}
	return "", nil
}

func (l *WalletLogic) GetAllTransactions(req *types.TransactionReq) (*pages.PageResult, error) {
	// rpc调用
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)

	in := &wallet.AssetReq{
		UserId:    userId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageNo:    req.PageNo,
		PageSize:  req.PageSize,
		Type:      req.Type,
		Symbol:    req.Symbol,
	}
	resp, err := l.svcCtx.UCWalletRpc.GetAllTransactions(ctx, in)
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", in)
	}

	b := make([]any, len(resp.List))
	for i, v := range resp.List {
		b[i] = v
	}
	pageResult := pages.New(b, req.PageNo, req.PageSize, resp.Total)
	return pageResult, nil
}
