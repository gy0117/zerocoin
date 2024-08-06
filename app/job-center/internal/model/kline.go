package model

import (
	"zero-common/tools"
)

// Kline 存入mongo db
type Kline struct {
	Period       string  `bson:"period,omitempty"`
	OpenPrice    float64 `bson:"openPrice,omitempty"`    // 开盘价格
	HighestPrice float64 `bson:"highestPrice,omitempty"` // 最高价格
	LowestPrice  float64 `bson:"lowestPrice,omitempty"`  // 最低价格
	ClosePrice   float64 `bson:"closePrice,omitempty"`   // 收盘价格

	// 交易量，以张为单位；如果是衍生品合约，数值为合约的张数；如果是币币/币币杠杆，数值为交易货币的数量
	Vol float64 `bson:"vol,omitempty"`
	// 交易量，以币为单位；如果是衍生品合约，数值为交易货币的数量；如果是币币/币币杠杆，数值为计价货币的数量
	VolCcy float64 `bson:"volCcy,omitempty"`
	// 交易量，以计价货币为单位；如BTC-USDT和BTC-USDT-SWAP，单位均是USDT；BTC-USD-SWAP单位是USD
	VolCcyQuote float64 `bson:"volCcyQuote,omitempty"`
	Time        int64   `bson:"time,omitempty"`    // 开始时间，Unix时间戳，毫秒
	TimeStr     string  `bson:"timeStr,omitempty"` // 时间，格式化
}

func NewKline(data []string, period string) *Kline {
	ts := tools.ToInt64(data[0])
	return &Kline{
		Period:       period,
		OpenPrice:    tools.ToFloat64(data[1]),
		HighestPrice: tools.ToFloat64(data[2]),
		LowestPrice:  tools.ToFloat64(data[3]),
		ClosePrice:   tools.ToFloat64(data[4]),
		Vol:          tools.ToFloat64(data[5]),
		VolCcy:       tools.ToFloat64(data[6]),
		VolCcyQuote:  tools.ToFloat64(data[7]),
		Time:         ts,
		TimeStr:      tools.ToTimeString(ts),
	}
}

func (k *Kline) Table(symbol, period string) string {
	return "exchange_kline_" + symbol + "_" + period
}
