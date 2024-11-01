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

func (od *OrderDaoV2) UpdateStatus(ctx context.Context, tx *gorm.DB, uid int64, orderId string, status int) error {
	err := tx.WithContext(ctx).Model(&model.ExchangeOrder{}).Where("order_id=?", orderId).Update("status", status).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}
