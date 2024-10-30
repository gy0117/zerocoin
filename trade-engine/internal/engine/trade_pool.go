package engine

import (
	"trade-engine/internal/model"
	"zero-common/kafka"
)

// TradePools 撮合交易池
type TradePools struct {
	pools map[string]*OrderBook
}

func (tp *TradePools) AddOrder(order *model.Order) error {
	orderBook, ok := tp.pools[order.TradePair]
	if !ok {
		return TradePairError
	}
	return orderBook.Add(order)
}

func (tp *TradePools) CancelOrder(tradePair string, id string) error {
	orderBook, ok := tp.pools[tradePair]
	if !ok {
		return TradePairError
	}
	return orderBook.Cancel(id)
}

func NewTradePools(tradePairs []string, kCli *kafka.KafkaClient) (*TradePools, error) {
	tp := TradePools{
		pools: make(map[string]*OrderBook),
	}

	for _, pair := range tradePairs {
		orderBook, err := NewOrderBook(pair, kCli)
		if err != nil {
			return nil, err
		}
		tp.pools[pair] = orderBook
		go orderBook.Start()
	}
	return &tp, nil
}
