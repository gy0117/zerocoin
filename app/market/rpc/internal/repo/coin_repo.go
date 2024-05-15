package repo

import (
	"context"
	"market-rpc/internal/model"
)

type CoinRepo interface {
	FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error)
	FindCoinByCoinId(ctx context.Context, coinId int64) (*model.Coin, error)
}
