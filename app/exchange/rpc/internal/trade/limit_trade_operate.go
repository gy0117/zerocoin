package trade

import (
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/trade/queue"
	"sort"
)

type LimitTradeOperate struct {
	// 这两个顺序的规定，都是为了快速匹配进行买卖
	// LimitPriceQueue，按照价格降序排列
	buyLimitQueue  *queue.LimitPriceQueue // 以限价买入队列，根据价格从高到低
	sellLimitQueue *queue.LimitPriceQueue // 以限价卖出队列，根据价格从低到高，这个在后面排序的时候，会处理的
}

func NewLimitTradeOperate() *LimitTradeOperate {
	return &LimitTradeOperate{
		buyLimitQueue:  &queue.LimitPriceQueue{},
		sellLimitQueue: &queue.LimitPriceQueue{},
	}
}

// JoinInBuyLimitQueue
// 将订单加入到 以限价买入队列
// 如果存在相同价格，则加入到对应的LimitPriceMap中；否则新建一个LimitPriceMap，并将其加入到queue中
// price - map
func (ltOperate *LimitTradeOperate) JoinInBuyLimitQueue(o *model.ExchangeOrder) {
	flag := false
	// 遍历限价队列，根据当前订单的价格，找到对应的位置，并放入
	for _, v := range ltOperate.buyLimitQueue.List {
		if o.Price == v.Price {
			v.List = append(v.List, o)
			flag = true
			break
		}
	}
	if !flag {
		lpm := &queue.LimitPriceMap{
			Price: o.Price,
			List:  []*model.ExchangeOrder{o},
		}
		ltOperate.buyLimitQueue.List = append(ltOperate.buyLimitQueue.List, lpm)
	}
}

// JoinInSellLimitQueue
// 将订单加入到 以限价卖出队列
// 如果存在相同价格，则加入到对应的LimitPriceMap中；否则新建一个LimitPriceMap，并将其加入到queue中
// price - map
func (ltOperate *LimitTradeOperate) JoinInSellLimitQueue(o *model.ExchangeOrder) {
	flag := false
	// 遍历限价队列，根据当前订单的价格，找到对应的位置，并放入
	for _, v := range ltOperate.sellLimitQueue.List {
		if o.Price == v.Price {
			v.List = append(v.List, o)
			flag = true
			break
		}
	}
	if !flag {
		lpm := &queue.LimitPriceMap{
			Price: o.Price,
			List:  []*model.ExchangeOrder{o},
		}
		ltOperate.sellLimitQueue.List = append(ltOperate.sellLimitQueue.List, lpm)
	}
}

func (ltOperate *LimitTradeOperate) Sort() {
	// 限价买队列，从高到低排序
	sort.Sort(ltOperate.buyLimitQueue.List)

	// 限价卖队列，从低到高排序
	sort.Sort(sort.Reverse(ltOperate.sellLimitQueue.List))
}
