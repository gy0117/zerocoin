package repo

import (
	"context"
	"exchange-rpc/internal/model"
	"zero-common/zerodb"
)

type OrderRepo interface {
	GetHistoryOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) (list []*model.ExchangeOrder, total int64, err error)
	GetCurrentOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) (list []*model.ExchangeOrder, total int64, err error)
	FindCurrentTradingCount(ctx context.Context, memId int64, symbol string, direction int) (int64, error)
	Save(ctx context.Context, conn zerodb.DbConn, order *model.ExchangeOrder) error
	FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error)
	CancelOrder(ctx context.Context, orderId string) error
	UpdateOrderStatus(ctx context.Context, orderId string, status int) error
	FindTradingOrderBySymbol(ctx context.Context, symbol string) ([]*model.ExchangeOrder, error)
	UpdateOrderComplete(ctx context.Context, orderId string, tradedAmount float64, turnover float64) error
}
