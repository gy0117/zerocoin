package model

import (
	"github.com/shopspring/decimal"
	"trade-engine/internal/queue"
)

type Side int

const (
	Buy = iota
	Sell
)

const (
	SellStr = "sell"
	BuyStr  = "buy"
)

func (s Side) String() string {
	if s == Buy {
		return BuyStr
	}
	return SellStr
}

type Type int

const (
	CancelOrder         = iota // 撤单
	MarketOrder                // 市价单
	LimitOrder                 // 限价单
	StopMarketOrder            // 市价触发单
	StopLimitOrder             // 限价触发单
	StopMarketLossOrder        // 市价止损单
	StopLimitLossOrder         // 限价止损单
	TrailingStopOrder
)

const (
	CancelOrderStr         = "cancel"
	MarketOrderStr         = "market"
	LimitOrderStr          = "limit"
	StopMarketOrderStr     = "stopMarketOrder"
	StopLimitOrderStr      = "stopLimitOrder"
	StopMarketLossOrderStr = "stopMarketLossOrder"
	StopLimitLossOrderStr  = "stopMarketLossOrder"
)

func (t Type) String() string {
	res := CancelOrderStr
	switch t {
	case MarketOrder:
		res = MarketOrderStr
	case LimitOrder:
		res = LimitOrderStr
	case StopMarketOrder:
		res = StopMarketOrderStr
	case StopLimitOrder:
		res = StopLimitOrderStr
	case StopMarketLossOrder:
		res = StopMarketLossOrderStr
	case StopLimitLossOrder:
		res = StopLimitLossOrderStr
	}
	return res
}

type Order struct {
	Id        string          `json:"id"`
	Uid       int64           `json:"uid"`
	TradePair string          `json:"tradePair"` // 交易对
	Price     decimal.Decimal `json:"price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Side      Side            `json:"side"` // 订单方向 buy or sell
	Type      Type            `json:"type"` // 订单类型
}

func (o *Order) SetQuantity(quantity decimal.Decimal) {
	o.Quantity = quantity
}

func (o *Order) GetUid() int64 {
	return o.Uid
}

func (o *Order) GetId() string {
	return o.Id
}

func (o *Order) GetQuantity() decimal.Decimal {
	return o.Quantity
}

var _ queue.NodeValue = (*Order)(nil)
