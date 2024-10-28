package logic

import (
	"context"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/svc"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
)

type Send2PlateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	kafkaDomain *domain.KafkaDomain
}

func NewSend2PlateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Send2PlateLogic {
	return &Send2PlateLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		kafkaDomain: domain.NewKafkaDomain(svcCtx.KCli, domain.NewOrderDomain(svcCtx.DB)),
	}
}

func (l *Send2PlateLogic) Send2Plate(req *order.SendOrderRequest) (*order.Empty, error) {
	logx.Infof("saga -> 发送到买卖盘, orderId: %s", req.OrderId)
	err := l.kafkaDomain.Send2Plate(req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("exchange-rpc.Send2Plate | 发送到kafka失败, err: %v", err)
	}
	return &order.Empty{}, nil
}

// TODO
func (l *Send2PlateLogic) Send2PlateRevert(req *order.SendOrderRequest) (*order.Empty, error) {
	logx.Infof("saga -> 撤销 发送到买卖盘, orderId: %s", req.OrderId)
	return &order.Empty{}, nil
}
