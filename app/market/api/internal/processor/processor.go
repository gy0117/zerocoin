package processor

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"market-api/internal/model"
	"market-api/internal/processor/handler"
	"zero-common/kafka"
)

const KLINE_1M = "kline_1m"

type ProcessType int

const (
	kline ProcessType = iota
	tradePlate
)

const topicExchangeOrderTradePlate = "exchange_order_trade_plate"

type ProcessData struct {
	Type ProcessType
	Key  []byte
	Data []byte
}

type Processor interface {
	Process(data ProcessData)
	AddHandler(h handler.MarketHandler)
	GetThumb() any
}

type DefaultProcessor struct {
	handlers []handler.MarketHandler
	kafkaCli *kafka.KafkaClient
	thumbMap map[string]*market.CoinThumb
}

func NewDefaultProcessor(kafkaCli *kafka.KafkaClient) *DefaultProcessor {
	return &DefaultProcessor{
		handlers: make([]handler.MarketHandler, 0),
		kafkaCli: kafkaCli,
		thumbMap: make(map[string]*market.CoinThumb),
	}
}

func (p *DefaultProcessor) Process(data ProcessData) {
	if data.Type == kline {
		// 交给
		symbol := string(data.Key)

		mk := &model.Kline{}
		_ = json.Unmarshal(data.Data, mk)

		for _, h := range p.handlers {
			h.HandleKline(symbol, mk, p.thumbMap)
		}
	} else if data.Type == tradePlate {
		logx.Info("收到买卖盘交易信息，准备利用socket.io推送到前端")
		symbol := string(data.Key)
		tradePlateData := &model.TradePlateResult{}
		_ = json.Unmarshal(data.Data, tradePlateData)
		for _, h := range p.handlers {
			h.HandleTradePlate(symbol, tradePlateData)
		}
	}
}

func (p *DefaultProcessor) AddHandler(h handler.MarketHandler) {
	p.handlers = append(p.handlers, h)
}

func (p *DefaultProcessor) Init(marketRpc mclient.Market) {
	p.startReadDataFromKafka(KLINE_1M, kline)
	p.startReadTradePlate(topicExchangeOrderTradePlate, tradePlate)
	p.initThumbMap(marketRpc)
}

// 读取买卖盘数据
func (p *DefaultProcessor) startReadTradePlate(topic string, processType ProcessType) {
	client := p.kafkaCli.StartRead(topic)
	go p.handleQueueData(client, processType)
}

func (p *DefaultProcessor) startReadDataFromKafka(topic string, pType ProcessType) {
	// 先start，再read
	client := p.kafkaCli.StartRead(topic)

	go p.handleQueueData(client, pType)
}

func (p *DefaultProcessor) handleQueueData(client *kafka.KafkaClient, pType ProcessType) {
	for {
		kfData := client.Read()
		pData := ProcessData{
			Type: pType,
			Key:  kfData.Key,
			Data: kfData.Data,
		}
		p.Process(pData)
	}
}

func (p *DefaultProcessor) GetThumb() any {
	cs := make([]*market.CoinThumb, len(p.thumbMap))
	i := 0
	for _, v := range p.thumbMap {
		cs[i] = v
		i++
	}
	return cs
}

func (p *DefaultProcessor) initThumbMap(marketRpc mclient.Market) {
	symbolThumbRes, err := marketRpc.FindCoinThumbTrend(context.Background(),
		&market.MarketRequest{})
	if err != nil {
		logx.Info(err)
	} else {
		coinThumbs := symbolThumbRes.List
		for _, v := range coinThumbs {
			p.thumbMap[v.Symbol] = v
		}
	}
}
