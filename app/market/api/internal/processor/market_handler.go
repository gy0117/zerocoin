package processor

import (
	"grpc-common/market/types/market"
	"market-api/internal/model"
)

type MarketHandler interface {
	HandleKline(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb)
	HandleTradePlate(symbol string, data *model.TradePlateResult)
}
