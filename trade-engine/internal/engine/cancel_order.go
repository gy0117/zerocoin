package engine

import (
	"github.com/shopspring/decimal"
	"time"
	"trade-engine/internal/model"
)

// TODO 加锁？
func (orderBook *OrderBook) cancel(orderId string) error {
	if score, ok := orderBook.mBid[orderId]; ok {
		return orderBook.cancelBid(score, orderId)
	}
	if score, ok := orderBook.mAsk[orderId]; ok {
		return orderBook.cancelAsk(score, orderId)
	}
	return CancelOrderError
}

func (orderBook *OrderBook) cancelBid(score decimal.Decimal, orderId string) error {
	node, _ := orderBook.bid.Find(score, orderId)
	if node == nil {
		return CancelOrderError
	}
	order, ok := node.GetValue().(*model.Order)
	if !ok {
		return CancelOrderError
	}
	orderBook.bid.Delete(score, orderId)
	trade := model.Trade{
		Id:             GenerateTradeId(),
		TradePair:      order.TradePair,
		MakerId:        order.Id,
		TakerId:        order.Id,
		MakerUser:      order.Uid,
		TakerUser:      order.Uid,
		Price:          order.Price.String(),
		Quantity:       order.Quantity.String(),
		TakerOrderSide: order.Side.String(),
		TakerOrderType: model.CancelOrderStr,
		Timestamp:      time.Now().UnixMilli(),
	}

	orderBook.PushTradeCanceled(trade)
	delete(orderBook.mBid, orderId)
	return nil
}

func (orderBook *OrderBook) cancelAsk(score decimal.Decimal, orderId string) error {
	node, _ := orderBook.ask.Find(score, orderId)
	if node == nil {
		return CancelOrderError
	}
	order, ok := node.GetValue().(*model.Order)
	if !ok {
		return CancelOrderError
	}
	orderBook.ask.Delete(score, orderId)
	trade := model.Trade{
		Id:             GenerateTradeId(),
		TradePair:      order.TradePair,
		MakerId:        order.Id,
		TakerId:        order.Id,
		MakerUser:      order.Uid,
		TakerUser:      order.Uid,
		Price:          order.Price.String(),
		Quantity:       order.Quantity.String(),
		TakerOrderSide: order.Side.String(),
		TakerOrderType: model.CancelOrderStr,
		Timestamp:      time.Now().UnixMilli(),
	}

	orderBook.PushTradeCanceled(trade)
	delete(orderBook.mAsk, orderId)
	return nil
}
