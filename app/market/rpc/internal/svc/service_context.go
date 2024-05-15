package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"market-rpc/internal/config"
	"market-rpc/internal/db"
	"zero-common/zerodb"
)

type ServiceContext struct {
	Config      config.Config
	Cache       cache.Cache
	DB          *zerodb.ZeroDB
	MongoClient *db.MongoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdsCache := cache.New(c.CacheRedis, nil, cache.NewStat("zerocoin"), nil, func(o *cache.Options) {})
	return &ServiceContext{
		Config:      c,
		Cache:       rdsCache,
		DB:          db.ConnMysql(c.Mysql.DataSource),
		MongoClient: db.ConnectMongo(c.Mongo),
	}
}
