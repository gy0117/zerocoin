package logic

import (
	"context"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/svc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"time"
	"zero-common/zerodb/tran"
)

type GetOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	orderDomain *domain.OrderDomain
	transaction tran.Transaction
	kafkaDomain *domain.KafkaDomain
}

func NewGetOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrdersLogic {
	orderDomain := domain.NewOrderDomain(svcCtx.DB)

	return &GetOrdersLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		orderDomain: orderDomain,
		transaction: tran.NewTransaction(svcCtx.DB.Conn),
		kafkaDomain: domain.NewKafkaDomain(svcCtx.KCli, orderDomain),
	}
}

func (l *GetOrdersLogic) GetHistoryOrder(in *order.OrderReq) (*order.OrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	exchangeOrders, total, err := l.orderDomain.GetHistoryOrder(ctx, in.GetSymbol(), in.GetUserId(), in.GetPage(), in.GetPageSize())
	if err != nil {
		return nil, errors.Wrapf(ErrGetHistoryOrder, "exchange-rpc getHistoryOrder, uid: %d, symbol: %s", in.GetUserId(), in.GetSymbol())
	}

	list := make([]*order.ExchangeOrder, len(exchangeOrders))
	for i, v := range exchangeOrders {
		list[i] = &order.ExchangeOrder{
			Id:            v.Id,
			OrderId:       v.OrderId,
			Amount:        v.Amount,
			BaseSymbol:    v.BaseSymbol,
			CanceledTime:  v.CanceledTime,
			CoinSymbol:    v.CoinSymbol,
			CompletedTime: v.CompletedTime,
			Direction:     model.DirectionMap.Value(v.Direction),
			UserId:        v.UserId,
			Price:         v.Price,
			Status:        int32(v.Status),
			Symbol:        v.Symbol,
			Time:          v.Time,
			TradedAmount:  v.TradedAmount,
			Turnover:      v.Turnover,
			Type:          model.TypeMap.Value(v.Type),
			UseDiscount:   v.UseDiscount,
		}
	}

	resp := &order.OrderResp{
		List:  list,
		Total: total,
	}
	return resp, nil
}

func (l *GetOrdersLogic) GetCurrentOrder(in *order.OrderReq) (*order.OrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	exchangeOrders, total, err := l.orderDomain.GetCurrentOrder(ctx, in.GetSymbol(), in.GetUserId(), in.GetPage(), in.GetPageSize())
	if err != nil {
		return nil, errors.Wrapf(ErrGetCurrentOrder, "exchange-rpc uid: %d, symbol: %s", in.GetUserId(), in.GetSymbol())
	}

	list := make([]*order.ExchangeOrder, len(exchangeOrders))
	for i, v := range exchangeOrders {
		list[i] = &order.ExchangeOrder{
			Id:            v.Id,
			OrderId:       v.OrderId,
			Amount:        v.Amount,
			BaseSymbol:    v.BaseSymbol,
			CanceledTime:  v.CanceledTime,
			CoinSymbol:    v.CoinSymbol,
			CompletedTime: v.CompletedTime,
			Direction:     model.DirectionMap.Value(v.Direction),
			UserId:        v.UserId,
			Price:         v.Price,
			Status:        int32(v.Status),
			Symbol:        v.Symbol,
			Time:          v.Time,
			TradedAmount:  v.TradedAmount,
			Turnover:      v.Turnover,
			Type:          model.TypeMap.Value(v.Type),
			UseDiscount:   v.UseDiscount,
		}
	}

	resp := &order.OrderResp{
		List:  list,
		Total: total,
	}
	return resp, nil
}
