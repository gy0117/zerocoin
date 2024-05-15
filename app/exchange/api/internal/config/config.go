package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ExchangeRpc zrpc.RpcClientConf
	Jwt         JwtConf
}

type JwtConf struct {
	AccessSecret string
	AccessExpire int64
	Issuer       string
}
