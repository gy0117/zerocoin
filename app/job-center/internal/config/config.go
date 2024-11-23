package config

import (
	"common/kafka"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"job-center/internal/db"
)

type Config struct {
	Okx        OkxConfig
	Mongo      db.MongoConfig
	Kafka      kafka.KafkaConfig
	Redis      cache.CacheConf
	Bitcoin    BitcoinConfig
	UCenterRpc zrpc.RpcClientConf
	LogConfig  logx.LogConf
}

type OkxConfig struct {
	ApiKey     string
	SecretKey  string
	Passphrase string
	Host       string
	Proxy      string
}

type BitcoinConfig struct {
	Url string
}
