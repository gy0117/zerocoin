package plate

import (
	"exchange-rpc/internal/model"
	"sync"
	"zero-common/operate"
)

const (
	BUY = iota
	SELL
)

// TradePlateResult 展示到买卖盘上的数据
type TradePlateResult struct {
	Direction    string            `json:"direction"`
	MaxAmount    float64           `json:"maxAmount"`
	MinAmount    float64           `json:"minAmount"`
	HighestPrice float64           `json:"highestPrice"`
	LowestPrice  float64           `json:"lowestPrice"`
	Symbol       string            `json:"symbol"`
	Items        []*TradePlateItem `json:"items"`
}

type TradePlateItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// TradePlate 盘口信息
type TradePlate struct {
	Items     []*TradePlateItem `json:"items"`
	Symbol    string
	direction int
	maxDepth  int // 最多有多少待交易信息
	lock      sync.RWMutex
}

func (p *TradePlate) Add(order *model.ExchangeOrder) {
	// 遍历盘口数据，将价格一致的添加到一起
	if p.direction != order.Direction {
		return
	}

	if order.Type == model.MarketPriceInt {
		// 市价 不进入买卖盘，市价单是以市场当前价格立即买入或卖出资产的订单
		// 限价单则是达到特定或更好价格后再执行的订单

		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, item := range p.Items {
		if item.Price == order.Price {
			// order.Amount委托的数量
			// order.TradedAmount 已经成交的数量
			item.Amount = operate.FloorFloat(item.Amount+order.Amount-order.TradedAmount, 8)
			return
		}
	}

	// 如果没有相同价格的盘口item，则新建一个
	size := len(p.Items)
	if size < p.maxDepth {
		item := &TradePlateItem{
			Amount: operate.FloorFloat(order.Amount-order.TradedAmount, 8),
			Price:  order.Price,
		}
		p.Items = append(p.Items, item)
	}
}

func (p *TradePlate) NewResult() *TradePlateResult {
	result := &TradePlateResult{}
	if p.direction == BUY {
		result.Direction = "BUY"
	} else if p.direction == SELL {
		result.Direction = "SELL"
	}
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol

	result.Items = p.Items[:len(p.Items)]
	return result
}

func (p *TradePlate) getMaxAmount() float64 {
	n := len(p.Items)
	if n <= 0 {
		return 0
	}
	max := p.Items[0].Amount
	for i := 1; i < n; i++ {
		curAmount := p.Items[i].Amount
		if curAmount > max {
			max = curAmount
		}
	}
	return max
}

func (p *TradePlate) getMinAmount() float64 {
	n := len(p.Items)
	if n <= 0 {
		return 0
	}
	min := p.Items[0].Amount
	for i := 1; i < n; i++ {
		curAmount := p.Items[i].Amount
		if curAmount < min {
			min = curAmount
		}
	}
	return min
}

func (p *TradePlate) getHighestPrice() float64 {
	n := len(p.Items)
	if n <= 0 {
		return 0
	}
	high := p.Items[0].Price
	for i := 1; i < n; i++ {
		cur := p.Items[i].Price
		if cur > high {
			high = cur
		}
	}
	return high
}

func (p *TradePlate) getLowestPrice() float64 {
	n := len(p.Items)
	if n <= 0 {
		return 0
	}
	low := p.Items[0].Price
	for i := 1; i < n; i++ {
		cur := p.Items[i].Price
		if cur < low {
			low = cur
		}
	}
	return low
}

func (p *TradePlate) Remove(order *model.ExchangeOrder, amount float64) {
	for i, item := range p.Items {
		if item.Price == order.Price {
			item.Amount = operate.SubFloor(item.Amount, amount, 8)
			if item.Amount <= 0 {
				p.Items = append(p.Items[:i], p.Items[i+1:]...)
			}
			break
		}
	}
}

func NewTradePlate(symbol string, direction int) *TradePlate {
	tp := &TradePlate{
		Symbol:    symbol,
		direction: direction,
		maxDepth:  100,
	}
	return tp
}
