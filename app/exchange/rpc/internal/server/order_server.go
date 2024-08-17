// Code generated by goctl. DO NOT EDIT.
// Source: order.proto

package server

import (
	"context"
	"exchange-rpc/internal/logic"
	"exchange-rpc/internal/svc"
	"grpc-common/exchange/types/order"
)


type OrderServer struct {
	svcCtx *svc.ServiceContext
	order.UnimplementedOrderServiceServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
	}
}

func (o *OrderServer) GetHistoryOrder(ctx context.Context, req *order.OrderReq) (*order.OrderResp, error) {
	l := logic.NewGetOrdersLogic(ctx, o.svcCtx)
	return l.GetHistoryOrder(req)
}

func (o *OrderServer) GetCurrentOrder(ctx context.Context, req *order.OrderReq) (*order.OrderResp, error) {
	l := logic.NewGetOrdersLogic(ctx, o.svcCtx)
	return l.GetCurrentOrder(req)
}

func (o *OrderServer) AddOrder(ctx context.Context, req *order.OrderReq) (*order.AddOrderResp, error) {
	l := logic.NewOrderLogic(ctx, o.svcCtx)
	return l.AddOrder(req)
}

func(o *OrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.AddOrderResp, error) {
	l := logic.NewCreateOrderLogic(ctx, o.svcCtx)
	return l.CreateOrder(req)
}

func (o *OrderServer) CreateOrderRevert(ctx context.Context, req *order.CreateOrderRequest) (*order.AddOrderResp, error) {
	l := logic.NewCreateOrderLogic(ctx, o.svcCtx)
	return l.CreateOrderRevert(req)
}


func (o *OrderServer) FindByOrderId(ctx context.Context, req *order.OrderReq) (*order.ExchangeOrder, error) {
	l := logic.NewOrderLogic(ctx, o.svcCtx)
	return l.FindByOrderId(req)
}

func (o *OrderServer) CancelOrder(ctx context.Context, req *order.OrderReq) (*order.CancelOrderResp, error) {
	l := logic.NewOrderLogic(ctx, o.svcCtx)
	return l.CancelOrder(req)
}

func (o *OrderServer) SendOrder2Plate(ctx context.Context, req *order.SendOrderRequest) (*order.Empty, error) {
	l := logic.NewSend2PlateLogic(ctx, o.svcCtx)
	return l.Send2Plate(req)
}
func (o *OrderServer) SendOrder2PlateRevert(ctx context.Context, req *order.SendOrderRequest) (*order.Empty, error) {
	l := logic.NewSend2PlateLogic(ctx, o.svcCtx)
	return l.Send2PlateRevert(req)
}