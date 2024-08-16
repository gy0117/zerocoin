package dao

import (
	"context"
	"exchange-rpc/internal/model"
	"gorm.io/gorm"
)

//const tableName = "exchange_order"

type OrderDaoV2 struct {
}

func NewOrderDaoV2() *OrderDaoV2 {
	return &OrderDaoV2{}
}

func (od *OrderDaoV2) CreateOrder(ctx context.Context, tx *gorm.DB, newOrder *model.ExchangeOrder) error {
	return tx.WithContext(ctx).Save(&newOrder).Error
}
