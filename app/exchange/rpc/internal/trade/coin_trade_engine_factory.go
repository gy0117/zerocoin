package trade

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"sync"
	"time"
	"zero-common/kafka"
	"zero-common/zerodb"
)

// 撮合交易引擎

type CoinTradeEngineFactory struct {
	coinTradeMap map[string]*CoinTradeEngine
	lock         sync.RWMutex
}

func NewCoinTradeEngineFactory() *CoinTradeEngineFactory {
	return &CoinTradeEngineFactory{
		coinTradeMap: make(map[string]*CoinTradeEngine),
	}
}

// Init 读取exchange_coin表中的信息，symbol，一个symbol对应一个撮合交易引擎
// 在exchange服务启动时，执行
// 撮合交易引擎初始化
func (f *CoinTradeEngineFactory) Init(marketRpc mclient.Market, cli *kafka.KafkaClient, db *zerodb.ZeroDB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := marketRpc.FindExchangeCoinVisible(ctx, &market.MarketRequest{})
	if err != nil {
		logx.Error(err)
		return
	}
	if len(resp.List) > 0 {
		for _, coin := range resp.List {
			f.AddCoinTrade(coin.Symbol, NewCoinTrade(coin.Symbol, cli, db))
		}
	}
}

// AddCoinTrade 添加交易引擎
// 一个symbol对应一个引擎，例如：BTC/USDT 对应一个针对BTC和USDT的交易引擎
func (f *CoinTradeEngineFactory) AddCoinTrade(symbol string, coinTrade *CoinTradeEngine) {
	f.lock.RLock()
	_, ok := f.coinTradeMap[symbol]
	if ok {
		logx.Info("[exchange-rpc] repeated addition of matching trading engines.")
		return
	}
	f.lock.RUnlock()

	f.lock.Lock()
	defer f.lock.Unlock()
	f.coinTradeMap[symbol] = coinTrade
}

// GetCoinTrade 获取symbol对应的交易引擎
func (f *CoinTradeEngineFactory) GetCoinTrade(symbol string) (*CoinTradeEngine, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	ct, ok := f.coinTradeMap[symbol]
	return ct, ok
}
