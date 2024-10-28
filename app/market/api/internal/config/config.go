package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-common/kafka"
)

type Config struct {
	rest.RestConf
	MarketRpc zrpc.RpcClientConf
	Kafka     kafka.KafkaConfig
	LogConfig logx.LogConf
}
