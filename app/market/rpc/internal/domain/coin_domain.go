package domain

import (
	"context"
	"errors"
	"market-rpc/internal/dao"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
	"zero-common/zerodb"
)

type CoinDomain struct {
	coinRepo repo.CoinRepo
}

func NewCoinDomain(db *zerodb.ZeroDB) *CoinDomain {
	return &CoinDomain{
		coinRepo: dao.NewCoinDao(db),
	}
}

func (d *CoinDomain) FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error) {
	coin, err := d.coinRepo.FindCoinInfo(ctx, unit)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, errors.New("查询的货币不存在")
	}
	return coin, nil
}

func (d *CoinDomain) FindCoinByCoinId(ctx context.Context, coinId int64) (*model.Coin, error) {
	coin, err := d.coinRepo.FindCoinByCoinId(ctx, coinId)
	if err != nil {
		return nil, err
	}
	if coin == nil {
		return nil, errors.New("查询的货币不存在")
	}
	return coin, nil
}
