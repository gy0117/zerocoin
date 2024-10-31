package engine

import (
	"github.com/shopspring/decimal"
	"time"
	"trade-engine/internal/model"
)

// 挂单 限价单
func (orderBook *OrderBook) handleBidLimit(order *model.Order) error {
	trades := make([]model.Trade, 0)
	for orderBook.askQueue.First() != nil &&
		orderBook.askQueue.First().GetScore().LessThanOrEqual(order.Price) &&
		order.Quantity.GreaterThan(decimal.Zero) {

		firstNode := orderBook.askQueue.First()
		nodeValue := firstNode.GetValue()

		if nodeValue.GetQuantity().GreaterThanOrEqual(order.Quantity) {
			// 全部吃掉order
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
				orderBook.askQueue.Delete(firstNode.GetScore(), nodeValue.GetId())
			}
		} else {
			// 吃掉firstNode，order还剩余
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
			orderBook.askQueue.Delete(firstNode.GetScore(), nodeValue.GetId())
		}
	}

	if order.Quantity.GreaterThan(decimal.Zero) {
		orderBook.bidQueue.Insert(order.Price, order)
		orderBook.bidQueue.AddToMap(order.Id, order.Price)
	}

	if len(trades) > 0 {
		orderBook.PushTradeTickets(trades...)
	}
	return nil
}

// 市价单
func (orderBook *OrderBook) handleBidMarket(order *model.Order) error {
	trades := make([]model.Trade, 0)
	// 市价单不看对方价格，就要立即成交
	for orderBook.askQueue.First() != nil && order.Quantity.GreaterThan(decimal.Zero) {
		firstNode := orderBook.askQueue.First()
		nodeValue := firstNode.GetValue()

		if nodeValue.GetQuantity().GreaterThanOrEqual(order.Quantity) {
			// 吃掉order
			trade := model.Trade{
				Id:             GenerateTradeId(),
				TradePair:      order.TradePair,
				MakerId:        nodeValue.GetId(), // 挂单id
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
				orderBook.askQueue.Delete(firstNode.GetScore(), nodeValue.GetId())
			}
		} else {
			// 吃掉firstNode
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
			orderBook.askQueue.Delete(firstNode.GetScore(), nodeValue.GetId())
		}
	}

	// check order是否完成被吃掉
	// TODO 具体看策略，应该加到市价单队列中 或者 推送出去
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

		//orderBook.bid.Insert(order.Price, &order)
		//orderBook.mBid[order.Id] = order.Price

		go orderBook.PushTradeCanceled(trade)
	}
	if len(trades) > 0 {
		orderBook.PushTradeTickets(trades...)
	}
	return nil
}
