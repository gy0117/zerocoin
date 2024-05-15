package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"market-rpc/internal/db"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql      MysqlConfig
	CacheRedis cache.CacheConf
	Mongo      db.MongoConfig
}

type MysqlConfig struct {
	DataSource string
}
