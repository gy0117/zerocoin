package engine

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"time"
	model2 "trade-engine/internal/model"
	"trade-engine/internal/queue"
	"zero-common/kafka"
)

const maxOrderCap = 1000000

const topicTradeTickets = "trade_tickets"

type OrderBook struct {
	tradePair  string
	bid        *queue.SkipListDesc        //  买入价，出价。即买家愿意出的价格。出价越高，越容易成交，因此从高到低排序
	ask        *queue.SkipList            //  卖出价，卖家出的价格。按从低到高的顺序排列
	mBid       map[string]decimal.Decimal // bid的订单id对应的score
	mAsk       map[string]decimal.Decimal // ask的订单id对应的score
	chanAdd    chan model2.Order          // 异步处理挂单逻辑
	chanCancel chan string                // 异步处理撤单逻辑

	kCli *kafka.KafkaClient
}

func NewOrderBook(tradePair string, kCli *kafka.KafkaClient) (*OrderBook, error) {
	bid, err := queue.NewSkipListDesc()
	if err != nil {
		return nil, err
	}
	ask, err := queue.NewSkipList()
	if err != nil {
		return nil, err
	}
	return &OrderBook{
		tradePair:  tradePair,
		bid:        bid,
		ask:        ask,
		mBid:       make(map[string]decimal.Decimal),
		mAsk:       make(map[string]decimal.Decimal),
		chanAdd:    make(chan model2.Order, maxOrderCap),
		chanCancel: make(chan string, maxOrderCap),
		kCli:       kCli,
	}, nil
}

func (orderBook *OrderBook) Start() {
	for {
		select {
		case order := <-orderBook.chanAdd:
			if err := orderBook.add(order); err != nil {
				fmt.Printf("OrderBook Start, add err: %v\n", err)
			}
		case orderId := <-orderBook.chanCancel:
			if err := orderBook.cancel(orderId); err != nil {
				fmt.Printf("OrderBook Start, cancel err: %v\n", err)
			}
		}
	}
}

// Add 挂单 异步
func (orderBook *OrderBook) Add(order *model2.Order) error {
	fmt.Println("orderBook Add")
	select {
	case orderBook.chanAdd <- *order:
		return nil
	case <-time.After(time.Second):
		return TimeoutError
	}
}

// Cancel 撤单 异步
func (orderBook *OrderBook) Cancel(id string) error {
	fmt.Println("orderBook Cancel")
	select {
	case orderBook.chanCancel <- id:
		return nil
	case <-time.After(time.Second):
		return TimeoutError
	}
}

// PushTradeTickets 发送成交单
func (orderBook *OrderBook) PushTradeTickets(trades ...model2.Trade) {
	marshal, _ := json.Marshal(trades)
	kData := kafka.KafkaData{
		Topic: topicTradeTickets,
		Data:  marshal,
	}
	orderBook.kCli.Send(kData)
	log.Printf("trade-ticket: %+v\n", trades)
}
