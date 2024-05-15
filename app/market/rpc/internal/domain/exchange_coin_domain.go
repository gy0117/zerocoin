package domain

import (
	"context"
	"errors"
	"market-rpc/internal/dao"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
	"zero-common/zerodb"
)

type ExchangeCoinDomain struct {
	exchangeCoin repo.ExchangeCoinRepo
}

func NewExchangeCoinDomain(db *zerodb.ZeroDB) *ExchangeCoinDomain {
	return &ExchangeCoinDomain{
		exchangeCoin: dao.NewExchangeCoinDao(db),
	}
}

// FindCoinVisible 从数据库中找到可见的coin
func (d *ExchangeCoinDomain) FindCoinVisible(ctx context.Context) []*model.ExchangeCoin {
	coins, err := d.exchangeCoin.FindCoinVisible(ctx)
	if err != nil {
		return make([]*model.ExchangeCoin, 0)
	}
	return coins
}

func (d *ExchangeCoinDomain) FindBySymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	exchangeCoin, err := d.exchangeCoin.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if exchangeCoin == nil {
		return nil, errors.New("交易对不存在")
	}
	return exchangeCoin, nil
}
