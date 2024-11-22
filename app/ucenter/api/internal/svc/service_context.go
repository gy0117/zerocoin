package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/uclient"
	"ucenter-api/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	redis         *redis.Redis
	UCRegisterRpc uclient.Register
	UCLoginRpc    uclient.Login
	UCWalletRpc   uclient.Wallet
	UserRpc       uclient.User
	MarketRpc     mclient.Market
	UCWithdrawRpc uclient.Withdraw
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		redis:         redis.MustNewRedis(c.CacheRedis[0].RedisConf),
		UCRegisterRpc: uclient.NewRegister(zrpc.MustNewClient(c.UCenterRpc)),
		UCLoginRpc:    uclient.NewLogin(zrpc.MustNewClient(c.UCenterRpc)),
		UCWalletRpc:   uclient.NewWallet(zrpc.MustNewClient(c.UCenterRpc)),
		UserRpc:       uclient.NewUser(zrpc.MustNewClient(c.UCenterRpc)),
		MarketRpc:     mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
		UCWithdrawRpc: uclient.NewWithdraw(zrpc.MustNewClient(c.UCenterRpc)),
	}
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
