package repo

import (
	"context"
	"job-center/internal/model"
)

type KlineRepo interface {
	SaveBatch(ctx context.Context, data []*model.Kline, symbol, period string) error
	DeleteGtTime(ctx context.Context, time int64, symbol string, period string) error
}
