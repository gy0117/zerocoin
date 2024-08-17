package domain

import (
	"context"
	"exchange-rpc/internal/dao/v2"
	"exchange-rpc/internal/model"
	"gorm.io/gorm"
)

type OrderDomainV2 struct {
	orderDao *dao.OrderDaoV2
}

func NewOrderDomainV2() *OrderDomainV2 {
	return &OrderDomainV2{
		orderDao: dao.NewOrderDaoV2(),
	}
}

func (od *OrderDomainV2) CreateOrder(ctx context.Context, tx *gorm.DB, newOrder *model.ExchangeOrder) error {
	return od.orderDao.CreateOrder(ctx, tx, newOrder)
}

func (od *OrderDomainV2) CreateOrderRevert(ctx context.Context, tx *gorm.DB, uid int64, orderId string, status int) error {
	return od.orderDao.UpdateStatus(ctx, tx, uid, orderId, status)
}
