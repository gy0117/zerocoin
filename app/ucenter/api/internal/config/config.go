package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UCenterRpc zrpc.RpcClientConf
	Jwt        JwtConf
	MarketRpc  zrpc.RpcClientConf
	LogConfig  logx.LogConf
}

type JwtConf struct {
	AccessSecret string
	AccessExpire int64
	Issuer       string
}
