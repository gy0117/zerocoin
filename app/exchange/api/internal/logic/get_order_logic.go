package logic

import (
	"context"
	"exchange-api/internal/page"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"time"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetHistoryOrders(req *types.ExchangeReq) (*page.PageData, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// rpc调用
	// 参数、返回值、方法
	userId := ctx.Value("userId").(int64)

	orderResp, err := l.svcCtx.OrderRpc.GetHistoryOrder(ctx, &order.OrderReq{
		Ip:       req.Ip,
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   userId,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "exchange-api GetHistoryOrders, uid: %d, req: %+v", userId, req)
	}

	b := make([]any, len(orderResp.List))
	for i, v := range orderResp.List {
		b[i] = v
	}
	data := page.New(b, req.PageNo, req.PageSize, orderResp.Total)
	return data, nil
}

func (l *GetOrderLogic) GetCurrentOrders(req *types.ExchangeReq) (*page.PageData, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// rpc调用
	// 参数、返回值、方法
	userId := ctx.Value("userId").(int64)

	orderResp, err := l.svcCtx.OrderRpc.GetCurrentOrder(ctx, &order.OrderReq{
		Ip:       req.Ip,
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   userId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "exchange-api GetCurrentOrders, uid: %d, req: %+v", userId, req)
	}

	b := make([]any, len(orderResp.List))
	for i, v := range orderResp.List {
		b[i] = v
	}
	data := page.New(b, req.PageNo, req.PageSize, orderResp.Total)
	return data, nil
}
