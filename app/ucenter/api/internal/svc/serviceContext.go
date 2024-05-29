package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/uclient"
	"ucenter-api/internal/config"
)

type ServiceContext struct {
	Config        config.Config
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
		UCRegisterRpc: uclient.NewRegister(zrpc.MustNewClient(c.UCenterRpc)),
		UCLoginRpc:    uclient.NewLogin(zrpc.MustNewClient(c.UCenterRpc)),
		UCWalletRpc:   uclient.NewWallet(zrpc.MustNewClient(c.UCenterRpc)),
		UserRpc:       uclient.NewUser(zrpc.MustNewClient(c.UCenterRpc)),
		MarketRpc:     mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
		UCWithdrawRpc: uclient.NewWithdraw(zrpc.MustNewClient(c.UCenterRpc)),
	}
}
