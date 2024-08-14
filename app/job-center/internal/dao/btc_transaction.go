package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"job-center/internal/model"
	"job-center/internal/repo"
)

var _ repo.BtcTransactionRepo = (*BtcTransactionDao)(nil)

type BtcTransactionDao struct {
	db *mongo.Database
}

func NewBtcTransactionDao(db *mongo.Database) *BtcTransactionDao {
	return &BtcTransactionDao{
		db: db,
	}
}

func (dao *BtcTransactionDao) Save(bt *model.BtcTransaction) error {
	collection := dao.db.Collection(model.TableBtcTransaction)
	_, err := collection.InsertOne(context.Background(), &bt)
	return err
}

func (dao *BtcTransactionDao) GetByTxId(txId string) (*model.BtcTransaction, error) {
	collection := dao.db.Collection(model.TableBtcTransaction)
	filter := bson.D{{"txId", txId}}
	var bt model.BtcTransaction
	err := collection.FindOne(context.Background(), filter).Decode(&bt)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &bt, nil
}
