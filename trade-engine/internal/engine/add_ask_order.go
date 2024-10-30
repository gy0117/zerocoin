package engine

import (
	"github.com/shopspring/decimal"
	"time"
	"trade-engine/internal/model"
)

func (orderBook *OrderBook) handleAskLimit(order model.Order) error {
	trades := make([]model.Trade, 0)
	for orderBook.bid.First() != nil &&
		orderBook.bid.First().GetScore().GreaterThanOrEqual(order.Price) &&
		order.Quantity.GreaterThan(decimal.Zero) {

		firstNode := orderBook.bid.First()
		nodeValue := firstNode.GetValue()

		if nodeValue.GetQuantity().GreaterThanOrEqual(order.Quantity) {
			trade := model.Trade{
				Id:             GenerateTradeId(),
				TradePair:      order.TradePair,
				MakerId:        nodeValue.GetId(),
				TakerId:        order.Id,
				MakerUser:      nodeValue.GetUid(),
				TakerUser:      order.Uid,
				Price:          firstNode.GetScore().String(),
				Quantity:       order.Quantity.String(),
				TakerOrderSide: order.Side.String(),
				TakerOrderType: order.Type.String(),
				Timestamp:      time.Now().UnixMilli(),
			}
			trades = append(trades, trade)

			left := nodeValue.GetQuantity().Sub(order.Quantity)
			order.Quantity = order.Quantity.Sub(order.Quantity)
			if left.GreaterThan(decimal.Zero) {
				nodeValue.SetQuantity(left)
			} else {
				orderBook.bid.Delete(firstNode.GetScore(), nodeValue.GetId())
			}
		} else {
			trade := model.Trade{
				Id:             GenerateTradeId(),
				TradePair:      order.TradePair,
				MakerId:        nodeValue.GetId(),
				TakerId:        order.Id,
				MakerUser:      nodeValue.GetUid(),
				TakerUser:      order.Uid,
				Price:          firstNode.GetScore().String(),
				Quantity:       nodeValue.GetQuantity().String(),
				TakerOrderSide: order.Side.String(),
				TakerOrderType: order.Type.String(),
				Timestamp:      time.Now().UnixMilli(),
			}
			trades = append(trades, trade)
			order.Quantity = order.Quantity.Sub(nodeValue.GetQuantity())

			orderBook.bid.Delete(firstNode.GetScore(), nodeValue.GetId())
		}
	}

	if order.Quantity.GreaterThan(decimal.Zero) {
		orderBook.ask.Insert(order.Price, &order)
		orderBook.mAsk[order.Id] = order.Price
	}
	if len(trades) > 0 {
		orderBook.PushTradeTickets(trades...)
	}
	return nil
}

func (orderBook *OrderBook) handleAskMarket(order model.Order) error {
	trades := make([]model.Trade, 0)
	for orderBook.bid.First() != nil && order.Quantity.GreaterThan(decimal.Zero) {
		firstNode := orderBook.bid.First()
		nodeValue := firstNode.GetValue()

		if nodeValue.GetQuantity().GreaterThanOrEqual(order.Quantity) {
			trade := model.Trade{
				Id:             GenerateTradeId(),
				TradePair:      order.TradePair,
				MakerId:        nodeValue.GetId(),
				TakerId:        order.Id,
				MakerUser:      nodeValue.GetUid(),
				TakerUser:      order.Uid,
				Price:          firstNode.GetScore().String(),
				Quantity:       order.Quantity.String(),
				TakerOrderSide: order.Side.String(),
				TakerOrderType: order.Type.String(),
				Timestamp:      time.Now().UnixMilli(),
			}
			trades = append(trades, trade)

			left := nodeValue.GetQuantity().Sub(order.Quantity)
			order.Quantity = order.Quantity.Sub(order.Quantity)
			if left.GreaterThan(decimal.Zero) {
				nodeValue.SetQuantity(left)
			} else {
				orderBook.bid.Delete(firstNode.GetScore(), nodeValue.GetId())
			}
		} else {
			trade := model.Trade{
				Id:             GenerateTradeId(),
				TradePair:      order.TradePair,
				MakerId:        nodeValue.GetId(),
				TakerId:        order.Id,
				MakerUser:      nodeValue.GetUid(),
				TakerUser:      order.Uid,
				Price:          firstNode.GetScore().String(),
				Quantity:       nodeValue.GetQuantity().String(),
				TakerOrderSide: order.Side.String(),
				TakerOrderType: order.Type.String(),
				Timestamp:      time.Now().UnixMilli(),
			}
			trades = append(trades, trade)
			order.Quantity = order.Quantity.Sub(nodeValue.GetQuantity())

			orderBook.bid.Delete(firstNode.GetScore(), nodeValue.GetId())
		}
	}

	// 同bid market，具体看怎么处理
	// 我觉得应该加到市价单队列，单独搞一个队列
	if order.Quantity.GreaterThan(decimal.Zero) {
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
		trades = append(trades, trade)
	}

	if len(trades) > 0 {
		orderBook.PushTradeTickets(trades...)
	}
	return nil
}
