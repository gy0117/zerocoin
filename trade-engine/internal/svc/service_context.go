package svc

import (
	"common/kafka"
	"trade-engine/internal/config"
)

type ServiceContext struct {
	Config config.Config
	KCli   *kafka.KafkaClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	kfCli := kafka.NewKafkaClient(c.Kafka)

	svc := &ServiceContext{
		Config: c,
		KCli:   kfCli,
	}
	return svc
}
