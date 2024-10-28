package config

import (
	"exchange-rpc/internal/db"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-common/kafka"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql      MysqlConfig
	CacheRedis cache.CacheConf
	Mongo      db.MongoConfig
	UCenter    zrpc.RpcClientConf
	Market     zrpc.RpcClientConf
	Kafka      kafka.KafkaConfig
	LogConfig  logx.LogConf
	DtmConf    zrpc.RpcClientConf
}

type MysqlConfig struct {
	DataSource string
	Driver     string
	Host       string
	Port       int64
	User       string
	Password   string
	Db         string
}
