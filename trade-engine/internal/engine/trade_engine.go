package engine

import (
	"context"
	"sync"
	"trade-engine/internal/model"
	"zero-common/kafka"
)

// TradeEngine 撮合交易引擎
type TradeEngine struct {
	pools sync.Map
	kCli  *kafka.KafkaClient
}

func (engine *TradeEngine) load(pair string) *OrderBook {
	if value, ok := engine.pools.Load(pair); ok {
		orderBook, _ := value.(*OrderBook)
		return orderBook
	}
	orderBook, _ := NewOrderBook(pair, engine.kCli)
	engine.pools.Store(pair, orderBook)
	go func() {
		orderBook.Start()
	}()
	return orderBook
}

func (engine *TradeEngine) AddOrder(ctx context.Context, order *model.Order) error {
	orderBook := engine.load(order.TradePair)
	return orderBook.Add(ctx, order)
}

func (engine *TradeEngine) CancelOrder(ctx context.Context, tradePair string, id string) error {
	orderBook := engine.load(tradePair)
	return orderBook.Cancel(ctx, id)
}

func NewTradeEngine(tradePairs []string, kCli *kafka.KafkaClient) (*TradeEngine, error) {
	engine := &TradeEngine{
		kCli: kCli,
	}
	for _, pair := range tradePairs {
		orderBook, err := NewOrderBook(pair, kCli)
		if err != nil {
			return nil, err
		}
		engine.pools.Store(pair, orderBook)
		go orderBook.Start()
	}
	return engine, nil
}
