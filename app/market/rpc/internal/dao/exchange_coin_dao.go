package dao

import (
	"common/zerodb"
	"common/zerodb/gorms"
	"context"
	"gorm.io/gorm"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
)

var _ repo.ExchangeCoinRepo = (*ExchangeCoinDao)(nil)

type ExchangeCoinDao struct {
	conn *gorms.GormConn
}

func NewExchangeCoinDao(db *zerodb.ZeroDB) repo.ExchangeCoinRepo {
	return &ExchangeCoinDao{
		conn: gorms.New(db.Conn),
	}
}

func (dao *ExchangeCoinDao) FindCoinVisible(ctx context.Context) (list []*model.ExchangeCoin, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.ExchangeCoin{}).Where("visible=?", 1).Find(&list).Error
	return
}

func (dao *ExchangeCoinDao) FindBySymbol(ctx context.Context, symbol string) (ec *model.ExchangeCoin, err error) {
	session := dao.conn.Session(ctx)
	// 查exchange_coin表
	err = session.Model(&model.ExchangeCoin{}).Where("symbol=?", symbol).Take(&ec).Error
	// 处理特殊case，找不到记录
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
