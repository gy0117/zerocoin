package engine

import (
	"context"
	"fmt"
	"log"
	"trade-engine/internal/model"
	"trade-engine/internal/queue"
	"zero-common/kafka"
)

const maxOrderCap = 1000000

const topicTradeTickets = "trade_tickets"
const topicTradeCanceled = "trade_canceled"

type OrderBook struct {
	tradePair  string
	bidQueue   *queue.BidQueue
	askQueue   *queue.AskQueue
	orderChan  chan *model.Order // 异步处理挂单逻辑
	cancelChan chan string       // 异步处理撤单逻辑

	kCli *kafka.KafkaClient
}

func NewOrderBook(tradePair string, kCli *kafka.KafkaClient) (*OrderBook, error) {
	return &OrderBook{
		tradePair:  tradePair,
		bidQueue:   queue.NewBidQueue(),
		askQueue:   queue.NewAskQueue(),
		orderChan:  make(chan *model.Order, maxOrderCap),
		cancelChan: make(chan string, maxOrderCap),
		kCli:       kCli,
	}, nil
}

func (orderBook *OrderBook) Start() {
	for {
		select {
		case order := <-orderBook.orderChan:
			if err := orderBook.add(order); err != nil {
				fmt.Printf("OrderBook Start, add err: %v\n", err)
			}
		case orderId := <-orderBook.cancelChan:
			if err := orderBook.cancel(orderId); err != nil {
				fmt.Printf("OrderBook Start, cancel err: %v\n", err)
			}
		}
	}
}

// Add 挂单 异步
func (orderBook *OrderBook) Add(ctx context.Context, order *model.Order) error {
	select {
	case orderBook.orderChan <- order:
		return nil
	case <-ctx.Done():
		return TimeoutError
	}
}

// Cancel 撤单 异步
func (orderBook *OrderBook) Cancel(ctx context.Context, id string) error {
	select {
	case orderBook.cancelChan <- id:
		return nil
	case <-ctx.Done():
		return TimeoutError
	}
}

// PushTradeTickets 发送成交单
func (orderBook *OrderBook) PushTradeTickets(trades ...model.Trade) {
	//marshal, _ := json.Marshal(trades)
	//kData := kafka.KafkaData{
	//	Topic: topicTradeTickets,
	//	Data:  marshal,
	//}
	//orderBook.kCli.Send(kData)
	log.Printf("trade-ticket: %+v\n", trades)
}

func (orderBook *OrderBook) PushTradeCanceled(trade model.Trade) {
	//marshal, _ := json.Marshal(trade)
	//kData := kafka.KafkaData{
	//	Topic: topicTradeCanceled,
	//	Data:  marshal,
	//}
	//orderBook.kCli.Send(kData)
	log.Printf("trade-canceled: %+v\n", trade)
}
