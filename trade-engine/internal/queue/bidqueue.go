package queue

import "github.com/shopspring/decimal"

type BidQueue struct {
	bid  *SkipListDesc              //  买入价，出价。即买家愿意出的价格。出价越高，越容易成交，因此从高到低排序
	mBid map[string]decimal.Decimal // bid的订单id对应的score
}

func NewBidQueue() *BidQueue {
	bid, _ := NewSkipListDesc()
	return &BidQueue{
		bid:  bid,
		mBid: make(map[string]decimal.Decimal),
	}
}

func (bq *BidQueue) First() *SkipListNode {
	return bq.bid.First()
}

func (bq *BidQueue) GetScore(id string) (score decimal.Decimal, ok bool) {
	score, ok = bq.mBid[id]
	return
}

func (bq *BidQueue) Find(score decimal.Decimal, orderId string) (*SkipListNode, []*SkipListNode) {
	return bq.bid.Find(score, orderId)
}

func (bq *BidQueue) Insert(score decimal.Decimal, value NodeValue) {
	bq.bid.Insert(score, value)
}

func (bq *BidQueue) Delete(score decimal.Decimal, orderId string) {
	bq.bid.Delete(score, orderId)
}

func (bq *BidQueue) DeleteFromMap(orderId string) {
	delete(bq.mBid, orderId)
}

func (bq *BidQueue) AddToMap(orderId string, score decimal.Decimal) {
	bq.mBid[orderId] = score
}
