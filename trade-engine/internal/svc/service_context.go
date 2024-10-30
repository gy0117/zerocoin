package svc

import (
	"trade-engine/internal/config"
	"zero-common/kafka"
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
