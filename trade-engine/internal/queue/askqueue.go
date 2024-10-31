package queue

import "github.com/shopspring/decimal"

type AskQueue struct {
	ask  *SkipList                  //  卖出价，卖家出的价格。按从低到高的顺序排列
	mAsk map[string]decimal.Decimal // ask的订单id对应的score
}

func NewAskQueue() *AskQueue {
	ask, _ := NewSkipList()
	return &AskQueue{
		ask:  ask,
		mAsk: make(map[string]decimal.Decimal),
	}
}

func (aq *AskQueue) First() *SkipListNode {
	return aq.ask.First()
}

func (aq *AskQueue) GetScore(id string) (score decimal.Decimal, ok bool) {
	score, ok = aq.mAsk[id]
	return
}

func (aq *AskQueue) Find(score decimal.Decimal, orderId string) (*SkipListNode, []*SkipListNode) {
	return aq.ask.Find(score, orderId)
}

func (aq *AskQueue) Insert(score decimal.Decimal, value NodeValue) {
	aq.ask.Insert(score, value)
}

func (aq *AskQueue) Delete(score decimal.Decimal, orderId string) {
	aq.ask.Delete(score, orderId)
}

func (aq *AskQueue) DeleteFromMap(orderId string) {
	delete(aq.mAsk, orderId)
}

func (aq *AskQueue) AddToMap(orderId string, score decimal.Decimal) {
	aq.mAsk[orderId] = score
}
