package svc

import (
	"exchange-api/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/exchange/eclient"
)

type ServiceContext struct {
	Config   config.Config
	redis    *redis.Redis
	OrderRpc eclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	serviceCtx := &ServiceContext{
		Config:   c,
		redis:    redis.MustNewRedis(c.CacheRedis[0].RedisConf),
		OrderRpc: eclient.NewOrder(zrpc.MustNewClient(c.ExchangeRpc)),
	}
	return serviceCtx
}

// IsTokenInBlackList check token if in black list
func (serviceCtx *ServiceContext) IsTokenInBlackList(token string) bool {
	exists, err := serviceCtx.redis.Exists(token)
	if err != nil {
		logx.Error(err)
		return false
	}
	return exists
}
