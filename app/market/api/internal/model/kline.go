package model

import (
	"grpc-common/market/types/market"
	"zero-common/operate"
)

type Kline struct {
	Period       string  `bson:"period,omitempty" json:"period"`
	OpenPrice    float64 `bson:"openPrice,omitempty" json:"openPrice"`
	HighestPrice float64 `bson:"highestPrice,omitempty" json:"highestPrice"`
	LowestPrice  float64 `bson:"lowestPrice,omitempty" json:"lowestPrice"`
	ClosePrice   float64 `bson:"closePrice,omitempty" json:"closePrice"`
	Time         int64   `bson:"time,omitempty" json:"time"`
	Count        float64 `bson:"count,omitempty" json:"count"`       //成交笔数
	Volume       float64 `bson:"volume,omitempty" json:"volume"`     //成交量
	Turnover     float64 `bson:"turnover,omitempty" json:"turnover"` //成交额
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

func (k *Kline) InitCoinThumb(symbol string) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.High = k.HighestPrice
	ct.Volume = k.Volume
	ct.Turnover = k.Turnover
	ct.Low = k.LowestPrice
	ct.Zone = 0
	ct.Change = 0
	ct.Chg = 0.0
	ct.UsdRate = k.ClosePrice
	ct.BaseUsdRate = 1
	ct.Trend = make([]float64, 0)
	ct.DateTime = k.Time
	return ct
}

func (k *Kline) ToCoinThumb(symbol string, ct *market.CoinThumb) *market.CoinThumb {
	isSame := false
	if ct.Symbol == symbol && ct.DateTime == k.Time {
		//认为这是同一个数据
		isSame = true
	}
	if !isSame {
		ct.Close = k.ClosePrice
		ct.Open = k.OpenPrice
		if ct.High < k.HighestPrice {
			ct.High = k.HighestPrice
		}
		ct.Low = k.LowestPrice
		if ct.Low > k.LowestPrice {
			ct.Low = k.LowestPrice
		}
		ct.Zone = 0
		ct.Volume = operate.AddN(k.Volume, ct.Volume, 8)
		ct.Turnover = operate.AddN(k.Turnover, ct.Turnover, 8)
		ct.Change = k.LowestPrice - ct.Close
		ct.Chg = operate.MulN(operate.DivN(ct.Change, ct.Close, 5), 100, 3)
		ct.UsdRate = k.ClosePrice
		ct.BaseUsdRate = 1
		ct.Trend = append(ct.Trend, k.ClosePrice)
		ct.DateTime = k.Time
	}
	return ct
}
