package engine

import (
	"github.com/shopspring/decimal"
	"time"
	"trade-engine/internal/model"
)

func (orderBook *OrderBook) cancel(orderId string) error {
	if score, ok := orderBook.bidQueue.GetScore(orderId); ok {
		return orderBook.cancelBid(score, orderId)
	}
	if score, ok := orderBook.askQueue.GetScore(orderId); ok {
		return orderBook.cancelAsk(score, orderId)
	}
	return CancelOrderError
}

func (orderBook *OrderBook) cancelBid(score decimal.Decimal, orderId string) error {
	node, _ := orderBook.bidQueue.Find(score, orderId)
	if node == nil {
		return CancelOrderError
	}
	order, ok := node.GetValue().(*model.Order)
	if !ok {
		return CancelOrderError
	}
	orderBook.bidQueue.Delete(score, orderId)
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
	orderBook.bidQueue.DeleteFromMap(orderId)
	return nil
}

func (orderBook *OrderBook) cancelAsk(score decimal.Decimal, orderId string) error {
	node, _ := orderBook.askQueue.Find(score, orderId)
	if node == nil {
		return CancelOrderError
	}
	order, ok := node.GetValue().(*model.Order)
	if !ok {
		return CancelOrderError
	}
	orderBook.askQueue.Delete(score, orderId)
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
	orderBook.askQueue.DeleteFromMap(orderId)
	return nil
}
