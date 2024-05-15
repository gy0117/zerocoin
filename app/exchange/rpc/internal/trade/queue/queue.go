package queue

import (
	"exchange-rpc/internal/model"
	"sync"
)

// TradeTimeQueue 按照时间 排序
type TradeTimeQueue []*model.ExchangeOrder

func (t TradeTimeQueue) Len() int {
	return len(t)
}

func (t TradeTimeQueue) Less(i, j int) bool {
	// 按照时间升序
	return t[i].Time < t[j].Time
}

func (t TradeTimeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// LimitPriceQueue 限价队列
type LimitPriceQueue struct {
	lock sync.RWMutex
	List TradeQueue
}

func (q *LimitPriceQueue) Lock() {
	q.lock.Lock()
}

func (q *LimitPriceQueue) Unlock() {
	q.lock.Unlock()
}

// LimitPriceMap 当前price价格的所有订单
// 买入或者卖出订单，总会有相等价格的，比如：用户A以10u 买入btc，用户B也以10u买入btc
type LimitPriceMap struct {
	Price float64
	List  []*model.ExchangeOrder
}

type TradeQueue []*LimitPriceMap

func (t TradeQueue) Len() int {
	return len(t)
}

func (t TradeQueue) Less(i, j int) bool {
	// 价格降序
	return t[i].Price > t[j].Price
}

func (t TradeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
