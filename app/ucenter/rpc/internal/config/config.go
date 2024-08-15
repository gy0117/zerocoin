package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-common/kafka"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql         MysqlConfig
	CacheRedis    cache.CacheConf
	CaptchaVerify CaptchaVerifyConf
	Jwt           JwtConf
	MarketRpc     zrpc.RpcClientConf
	Kafka         kafka.KafkaConfig
	ExchangeRpc   zrpc.RpcClientConf
	LogConfig     logx.LogConf
}

type MysqlConfig struct {
	DataSource string
}

type CaptchaVerifyConf struct {
	Vid       string
	SecretKey string
}

type JwtConf struct {
	AccessSecret string
	AccessExpire int64
	Issuer       string
}
