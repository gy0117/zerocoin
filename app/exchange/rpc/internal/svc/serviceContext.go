package svc

import (
	"exchange-rpc/internal/config"
	"exchange-rpc/internal/db"
	"exchange-rpc/internal/kafka/consumer"
	"exchange-rpc/internal/trade"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/uclient"
	"zero-common/kafka"
	"zero-common/zerodb"
)

type ServiceContext struct {
	Config      config.Config
	Cache       cache.Cache
	DB          *zerodb.ZeroDB
	MongoClient *db.MongoClient
	UserRpc     uclient.User
	MarketRpc   mclient.Market
	WalletRpc   uclient.Wallet
	KCli        *kafka.KafkaClient
}

func (svc *ServiceContext) init() {
	factory := trade.NewCoinTradeEngineFactory()
	factory.Init(svc.MarketRpc, svc.KCli, svc.DB)

	rdsConf := svc.Config.CacheRedis[0].RedisConf
	redisCli := redis.MustNewRedis(rdsConf)

	kafkaConsumer := consumer.NewKafkaConsumer(svc.KCli, factory, svc.DB, redisCli)
	kafkaConsumer.Run()
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdsCache := cache.New(c.CacheRedis, nil, cache.NewStat("zerocoin"), nil, func(o *cache.Options) {})
	kfCli := kafka.NewKafkaClient(c.Kafka)

	svc := &ServiceContext{
		Config:      c,
		Cache:       rdsCache,
		DB:          db.ConnMysql(c.Mysql.DataSource),
		MongoClient: db.ConnectMongo(c.Mongo),
		UserRpc:     uclient.NewUser(zrpc.MustNewClient(c.UCenter)),
		MarketRpc:   mclient.NewMarket(zrpc.MustNewClient(c.Market)),
		WalletRpc:   uclient.NewWallet(zrpc.MustNewClient(c.UCenter)),
		KCli:        kfCli,
	}

	svc.init()

	return svc
}
