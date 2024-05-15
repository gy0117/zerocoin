package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/exchange/eclient"
	"grpc-common/market/mclient"
	"ucenter-rpc/internal/config"
	"ucenter-rpc/internal/db"
	"ucenter-rpc/internal/kafka/consumer"
	"zero-common/kafka"
	"zero-common/zerodb"
)

const (
	topicAddExchangeOrder               = "add_exchange_order" // 发布委托（创建订单）topic
	topicUpdateWalletAfterOrderComplete = "update_wallet_after_order_complete"
	topicBtcTransaction                 = "btc_transaction"
	topicBtcWithdraw                    = "btc_withdraw"
)

type ServiceContext struct {
	Config    config.Config
	Cache     cache.Cache
	DB        *zerodb.ZeroDB
	MarketRpc mclient.Market
	KCli      *kafka.KafkaClient
}

func (svc *ServiceContext) init() {
	rdsConf := svc.Config.CacheRedis[0].RedisConf
	newRedis := redis.MustNewRedis(rdsConf)

	cli := kafka.NewKafkaClient(svc.Config.Kafka)

	completeKCli := cli.StartRead(topicUpdateWalletAfterOrderComplete)
	go consumer.UpdateWalletAfterOrderComplete(completeKCli, newRedis, svc.DB)

	btCli := cli.StartRead(topicBtcTransaction)
	go consumer.ConsumeBtcTransaction(btCli, svc.DB)

	withdrawCli := cli.StartRead(topicBtcWithdraw)
	go consumer.ConsumeBtcWithdraw(withdrawCli, svc.DB)
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdsCache := cache.New(c.CacheRedis, nil, cache.NewStat("zerocoin"), nil, func(o *cache.Options) {})

	zeroDB := db.ConnMysql(c.Mysql.DataSource)
	orderRpc := eclient.NewOrder(zrpc.MustNewClient(c.ExchangeRpc))
	cli := kafka.NewKafkaClient(c.Kafka)
	kCli := cli.StartRead(topicAddExchangeOrder)

	rdsConf := c.CacheRedis[0].RedisConf
	newRedis := redis.MustNewRedis(rdsConf)

	go consumer.ExchangeOrderAdd(kCli, orderRpc, zeroDB, newRedis)

	s := &ServiceContext{
		Config:    c,
		Cache:     rdsCache,
		DB:        zeroDB,
		MarketRpc: mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
		KCli:      cli,
	}
	s.init()
	return s
}
