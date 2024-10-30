package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"zero-common/kafka"
)

type Config struct {
	zrpc.RpcServerConf
	LogConfig logx.LogConf
	Kafka     kafka.KafkaConfig
}
