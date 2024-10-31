package engine

import (
	"strconv"
	"time"
	"trade-engine/internal/model"
)

func (orderBook *OrderBook) add(order *model.Order) error {
	// 买入单
	if order.Side == model.Buy {
		return orderBook.handleBid(order)
	} else if order.Side == model.Sell {
		return orderBook.handleAsk(order)
	}
	return OrderSideError
}

func (orderBook *OrderBook) handleBid(order *model.Order) error {
	if order.Type == model.LimitOrder {
		return orderBook.handleBidLimit(order)
	} else if order.Type == model.MarketOrder {
		return orderBook.handleBidMarket(order)
	}
	return OrderTypeError
}

func (orderBook *OrderBook) handleAsk(order *model.Order) error {
	if order.Type == model.LimitOrder {
		return orderBook.handleAskLimit(order)
	} else if order.Type == model.MarketOrder {
		return orderBook.handleAskMarket(order)
	}
	return OrderTypeError
}

func GenerateTradeId() string {
	milli := time.Now().UnixMilli()
	return strconv.Itoa(int(milli))
}
