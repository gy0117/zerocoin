package model

import "zero-common/enum"

const (
	StatusTrading   = iota // 交易中
	StatusCompleted        // 完成
	StatusCanceled         // 取消
	StatusOverTimed        // 超时
	StatusInit
)

var StatusMap = enum.Enum{
	StatusTrading:   "TRADING",
	StatusCompleted: "COMPLETED",
	StatusCanceled:  "CANCELED",
	StatusOverTimed: "OVERTIMED",
}

const (
	LimitPrice     = "LIMIT_PRICE"
	MarketPrice    = "MARKET_PRICE"
	MarketPriceInt = 0
	LimitPriceInt  = 1
)

var TypeMap = enum.Enum{
	MarketPriceInt: MarketPrice,
	LimitPriceInt:  LimitPrice,
}

const (
	DirectionBuy     = "BUY"
	DirectionSell    = "SELL"
	DirectionBuyInt  = 0
	DirectionSellInt = 1
)

var DirectionMap = enum.Enum{
	DirectionBuyInt:  DirectionBuy,
	DirectionSellInt: DirectionSell,
}

type ExchangeOrder struct {
	Id            int64   `gorm:"column:id" json:"id"`
	OrderId       string  `gorm:"column:order_id" json:"orderId"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
	BaseSymbol    string  `gorm:"column:base_symbol" json:"baseSymbol"`
	CanceledTime  int64   `gorm:"column:canceled_time" json:"canceledTime"`
	CoinSymbol    string  `gorm:"column:coin_symbol" json:"coinSymbol"`
	CompletedTime int64   `gorm:"column:completed_time" json:"completedTime"`
	Direction     int     `gorm:"column:direction" json:"direction"`
	UserId        int64   `gorm:"column:user_id" json:"userId"`
	Price         float64 `gorm:"column:price" json:"price"`
	Status        int     `gorm:"column:status" json:"status"`
	Symbol        string  `gorm:"column:symbol" json:"symbol"`
	Time          int64   `gorm:"column:time" json:"time"`
	TradedAmount  float64 `gorm:"column:traded_amount" json:"tradedAmount"`
	Turnover      float64 `gorm:"column:turnover" json:"turnover"`
	Type          int     `gorm:"column:type" json:"type"`
	UseDiscount   string  `gorm:"column:use_discount" json:"useDiscount"`
}

func (*ExchangeOrder) TableName() string {
	return "exchange_order"
}

func NewOrder() *ExchangeOrder {
	return &ExchangeOrder{}
}

func TransferDirection(direction string) int {
	if direction == DirectionSell {
		return 1
	}
	return 0
}

func TransferType(tp string) int {
	if tp == LimitPrice {
		return 1
	}
	return 0
}
