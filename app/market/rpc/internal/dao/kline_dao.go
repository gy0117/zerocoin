package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
)

var _ repo.KlineRepo = (*KlineDao)(nil)

type KlineDao struct {
	mongoDb *mongo.Database
}

func NewKlineDao(db *mongo.Database) repo.KlineRepo {
	return &KlineDao{
		mongoDb: db,
	}
}

func (kd *KlineDao) FindBySymbolTime(ctx context.Context, symbol, period string, from, end int64, sort string) ([]*model.Kline, error) {
	mk := &model.Kline{}
	sortFlag := -1
	if sort == "asc" {
		sortFlag = 1
	}

	collection := kd.mongoDb.Collection(mk.Table(symbol, period))
	cur, err := collection.Find(
		ctx,
		bson.D{{"time", bson.D{{"$gte", from}, {"$lte", end}}}},
		&options.FindOptions{Sort: bson.D{{"time", sortFlag}}})
	if err != nil {
		return nil, err
	}
	var list []*model.Kline
	err = cur.All(ctx, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
