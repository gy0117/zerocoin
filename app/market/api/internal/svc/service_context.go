package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"market-api/internal/config"
	"market-api/internal/processor"
	"market-api/internal/ws"
	"zero-common/kafka"
)

type ServiceContext struct {
	Config          config.Config
	ExchangeRateRpc mclient.ExchangeRate
	MarketRpc       mclient.Market
	Processor       processor.Processor
}

func NewServiceContext(c config.Config, server *ws.WebSocketServer) *ServiceContext {

	market := mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc))

	// NewServiceContext只执行一次，因此在这里执行Processor的初始化
	kafkaCli := kafka.NewKafkaClient(c.Kafka)
	defaultProcessor := processor.NewDefaultProcessor(kafkaCli)
	defaultProcessor.Init(market)
	defaultProcessor.AddHandler(processor.NewWsHandler(server))

	return &ServiceContext{
		Config:          c,
		ExchangeRateRpc: mclient.NewExchangeRate(zrpc.MustNewClient(c.MarketRpc)),
		MarketRpc:       market,
		Processor:       defaultProcessor,
	}
}
