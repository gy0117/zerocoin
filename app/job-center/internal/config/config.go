package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"job-center/internal/db"
	"zero-common/kafka"
)

type Config struct {
	Okx        OkxConfig
	Mongo      db.MongoConfig
	Kafka      kafka.KafkaConfig
	Redis      cache.CacheConf
	Bitcoin    BitcoinConfig
	UCenterRpc zrpc.RpcClientConf
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
