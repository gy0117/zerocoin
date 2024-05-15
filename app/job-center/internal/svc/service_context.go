package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/ucenter/uclient"
	"job-center/internal/config"
	"job-center/internal/db"
	"zero-common/kafka"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *db.MongoClient
	KafkaClient *kafka.KafkaClient
	CacheRedis  cache.Cache
	AssetRpc    uclient.Wallet
}

func NewServiceContext(c config.Config) *ServiceContext {

	kfCli := kafka.NewKafkaClient(c.Kafka)
	kfCli.StartWrite()

	rdsCache := cache.New(c.Redis, nil, cache.NewStat("zerocoin"), nil, func(o *cache.Options) {})

	return &ServiceContext{
		Config:      c,
		MongoClient: db.ConnectMongo(c.Mongo),
		KafkaClient: kfCli,
		CacheRedis:  rdsCache,
		AssetRpc:    uclient.NewWallet(zrpc.MustNewClient(c.UCenterRpc)),
	}
}
