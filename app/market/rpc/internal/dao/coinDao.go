package dao

import (
	"context"
	"gorm.io/gorm"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
)

var _ repo.CoinRepo = (*CoinDao)(nil)

type CoinDao struct {
	conn *gorms.GormConn
}

func NewCoinDao(db *zerodb.ZeroDB) *CoinDao {
	return &CoinDao{
		conn: gorms.New(db.Conn),
	}
}

func (dao *CoinDao) FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error) {
	session := dao.conn.Session(ctx)
	coin := &model.Coin{}
	err := session.Model(&model.Coin{}).Where("unit=?", unit).Find(coin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return coin, err
}

func (dao *CoinDao) FindCoinByCoinId(ctx context.Context, coinId int64) (*model.Coin, error) {
	session := dao.conn.Session(ctx)
	coin := &model.Coin{}
	err := session.Model(&model.Coin{}).Where("id=?", coinId).Find(coin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return coin, err
}
