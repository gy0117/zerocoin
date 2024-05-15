package repo

import (
	"context"
	"market-rpc/internal/model"
)

type ExchangeCoinRepo interface {
	FindCoinVisible(ctx context.Context) ([]*model.ExchangeCoin, error)
	FindBySymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error)
}
