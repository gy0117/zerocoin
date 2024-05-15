package model

import (
	"grpc-common/market/types/market"
	"zero-common/operate"
)

type Kline struct {
	Period       string  `bson:"period,omitempty"`
	OpenPrice    float64 `bson:"openPrice,omitempty"`
	HighestPrice float64 `bson:"highestPrice,omitempty"`
	LowestPrice  float64 `bson:"lowestPrice,omitempty"`
	ClosePrice   float64 `bson:"closePrice,omitempty"`
	Time         int64   `bson:"time,omitempty"`
	Count        float64 `bson:"count,omitempty"`    //成交笔数
	Volume       float64 `bson:"volume,omitempty"`   //成交量
	Turnover     float64 `bson:"turnover,omitempty"` //成交额
}

func (k *Kline) Table(symbol, period string) string {
	return "exchange_kline_" + symbol + "_" + period
}

func (k *Kline) Kline2CoinThumb(symbol string, old *Kline) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.Zone = 0
	ct.Change = k.ClosePrice - old.ClosePrice
	ct.Chg = operate.MulN(operate.DivN(ct.Change, old.ClosePrice, 5), 100, 5)
	return ct
}

func DefaultCoinThumb(symbol string) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Trend = []float64{}
	return ct
}
