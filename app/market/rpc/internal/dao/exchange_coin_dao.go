package dao

import (
	"context"
	"gorm.io/gorm"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
	"zero-common/zerodb"
	"zero-common/zerodb/gorms"
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

func (dao *ExchangeCoinDao) FindBySymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	session := dao.conn.Session(ctx)
	ec := model.ExchangeCoin{}
	// 查exchange_coin表
	err := session.Model(&model.ExchangeCoin{}).Where("symbol=?", symbol).Find(&ec).Error
	// 处理特殊case，找不到记录
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ec, err
}
