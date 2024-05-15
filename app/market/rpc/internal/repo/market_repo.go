package repo

import (
	"context"
	"market-rpc/internal/model"
)

type KlineRepo interface {
	FindBySymbolTime(ctx context.Context, symbol, period string, from, end int64, sort string) ([]*model.Kline, error)
}
