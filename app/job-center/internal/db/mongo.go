package db

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoConfig struct {
	Url      string
	Username string
	Password string
	DbName   string
}

type MongoClient struct {
	cli *mongo.Client
	Db  *mongo.Database
}

// ConnectMongo mongo连接
func ConnectMongo(c MongoConfig) *MongoClient {
	mc := &MongoClient{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 认证
	credential := options.Credential{
		Username: c.Username,
		Password: c.Password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Url).SetAuth(credential))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	db := client.Database(c.DbName)

	mc.cli = client
	mc.Db = db
	return mc
}

func (c *MongoClient) DisConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := c.cli.Disconnect(ctx); err != nil {
		logx.Error(err)
	}
	fmt.Println("关注mongo连接...")
}
