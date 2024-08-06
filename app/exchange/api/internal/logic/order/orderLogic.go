package order

import (
	"context"
	"exchange-api/internal/page"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"time"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetHistoryOrders 分页
func (l *OrderLogic) GetHistoryOrders(req *types.ExchangeReq) (*page.PageData, error) {
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
		return nil, err
	}

	b := make([]any, len(orderResp.List))
	for i, v := range orderResp.List {
		b[i] = v
	}
	data := page.New(b, req.PageNo, req.PageSize, orderResp.Total)
	return data, nil
}

func (l *OrderLogic) GetCurrentOrders(req *types.ExchangeReq) (*page.PageData, error) {
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
		return nil, err
	}

	b := make([]any, len(orderResp.List))
	for i, v := range orderResp.List {
		b[i] = v
	}
	data := page.New(b, req.PageNo, req.PageSize, orderResp.Total)
	return data, nil
}

func (l *OrderLogic) AddOrder(req *types.ExchangeReq) (string, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)

	// rpc调用
	orderResp, err := l.svcCtx.OrderRpc.AddOrder(ctx, &order.OrderReq{
		Symbol:      req.Symbol,
		Price:       req.Price,
		Amount:      req.Amount,
		Direction:   req.Direction,
		Type:        req.Type,
		UseDiscount: int32(req.UseDiscount),
		UserId:      userId,
	})
	if err != nil {
		return "", err
	}
	return orderResp.OrderId, nil
}
