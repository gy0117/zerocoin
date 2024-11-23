package config

import (
	"common/kafka"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	LogConfig logx.LogConf
	Kafka     kafka.KafkaConfig
}
