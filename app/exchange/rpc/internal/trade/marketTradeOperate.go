package trade

import (
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/trade/queue"
	"sort"
	"sync"
)

// MarketTradeOperate 市价交易操作
type MarketTradeOperate struct {
	buyMarketQueue queue.TradeTimeQueue // 买入市价队列
	bmqLock        sync.RWMutex

	sellMarketQueue queue.TradeTimeQueue // 卖出市价队列
	smqLock         sync.RWMutex
}

func NewMarketTradeOperate() *MarketTradeOperate {
	return &MarketTradeOperate{}
}

func (mtOperate *MarketTradeOperate) JoinInBuyMarketQueueSafe(order *model.ExchangeOrder) {
	mtOperate.bmqLock.Lock()
	defer mtOperate.bmqLock.Unlock()

	mtOperate.buyMarketQueue = append(mtOperate.buyMarketQueue, order)
}

func (mtOperate *MarketTradeOperate) JoinInSellMarketQueueSafe(order *model.ExchangeOrder) {
	mtOperate.smqLock.Lock()
	defer mtOperate.smqLock.Unlock()

	mtOperate.sellMarketQueue = append(mtOperate.sellMarketQueue, order)
}

// Sort 按照时间排序
func (mtOperate *MarketTradeOperate) Sort() {
	sort.Sort(mtOperate.buyMarketQueue)
	sort.Sort(mtOperate.sellMarketQueue)
}

func (mtOperate *MarketTradeOperate) JoinInMarketPriceQueue(order *model.ExchangeOrder) {
	if order.Type != model.MarketPriceInt {
		return
	}
	if order.Direction == model.DirectionBuyInt {
		mtOperate.buyMarketQueue = append(mtOperate.buyMarketQueue, order)
		sort.Sort(mtOperate.buyMarketQueue)
	} else if order.Direction == model.DirectionSellInt {
		mtOperate.sellMarketQueue = append(mtOperate.sellMarketQueue, order)
		sort.Sort(mtOperate.sellMarketQueue)
	}
}
