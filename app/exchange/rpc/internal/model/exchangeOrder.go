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
	OrderId       string  `gorm:"column:order_id" json:"orderId"`             // 订单id
	Amount        float64 `gorm:"column:amount" json:"amount"`                // 买入或者卖出量
	BaseSymbol    string  `gorm:"column:base_symbol" json:"baseSymbol"`       // 结算单位
	CanceledTime  int64   `gorm:"column:canceled_time" json:"canceledTime"`   // 取消时间
	CoinSymbol    string  `gorm:"column:coin_symbol" json:"coinSymbol"`       // 币单位
	CompletedTime int64   `gorm:"column:completed_time" json:"completedTime"` // 完成时间
	Direction     int     `gorm:"column:direction" json:"direction"`          // 订单方向 0买 1卖
	UserId        int64   `gorm:"column:user_id" json:"userId"`
	Price         float64 `gorm:"column:price" json:"price"`                // 挂单价格
	Status        int     `gorm:"column:status" json:"status"`              // 订单状态 0 交易中 1 完成 2 取消 3 超时
	Symbol        string  `gorm:"column:symbol" json:"symbol"`              // 交易对
	Time          int64   `gorm:"column:time" json:"time"`                  // 挂单时间
	TradedAmount  float64 `gorm:"column:traded_amount" json:"tradedAmount"` // 成交量
	Turnover      float64 `gorm:"column:turnover" json:"turnover"`          // 成交额
	Type          int     `gorm:"column:type" json:"type"`                  // 挂单类型  0 市场价  1 最低价
	UseDiscount   string  `gorm:"column:use_discount" json:"useDiscount"`   // 是否使用折扣 0 不使用  1 使用
}

func (*ExchangeOrder) TableName() string {
	return "exchange_order"
}

func NewOrder() *ExchangeOrder {
	return &ExchangeOrder{}
}

func TransferDirection(direction string) int {
	//if direction == DirectionSell {
	//	return DirectionSellInt
	//}
	//return DirectionBuyInt

	for k, v := range DirectionMap {
		if v == direction {
			return k
		}
	}
	return -1
}

func TransferType(tp string) int {
	//if tp == LimitPrice {
	//	return LimitPriceInt
	//}
	//return MarketPriceInt

	for k, v := range TypeMap {
		if v == tp {
			return k
		}
	}
	return -1
}
