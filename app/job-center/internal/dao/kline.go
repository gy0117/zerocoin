package dao

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"job-center/internal/model"
	"strconv"
)

type KlineDao struct {
	db *mongo.Database
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}

func (d *KlineDao) SaveBatch(ctx context.Context, data []*model.Kline, symbol, period string) error {
	m := &model.Kline{}
	collection := d.db.Collection(m.Table(symbol, period))

	docs := make([]interface{}, len(data))

	for i, v := range data {
		docs[i] = v
	}

	_, err := collection.InsertMany(ctx, docs, nil)

	return err
}

func (d *KlineDao) DeleteGtTime(ctx context.Context, time int64, symbol string, period string) error {
	m := &model.Kline{}
	table := m.Table(symbol, period)
	collection := d.db.Collection(table)

	// $gte 大于等于
	filter := bson.D{{"time", bson.D{{"$gte", time}}}}
	deleteResult, err := collection.DeleteMany(ctx, filter)
	if err == nil {
		logx.Info("删除 " + table + " 表中的 " + strconv.FormatInt(deleteResult.DeletedCount, 10) + " 个数据")
	}
	return err
}
