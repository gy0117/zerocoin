package domain

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"job-center/internal/dao"
	"job-center/internal/db"
	"job-center/internal/model"
	"job-center/internal/repo"
	"time"
)

type KlineDomain struct {
	klineRepo repo.KlineRepo
}

func NewKlineDomain(client *db.MongoClient) *KlineDomain {
	return &KlineDomain{
		klineRepo: dao.NewKlineDao(client.Db),
	}
}

// SaveBatch 在domain中处理数据，dao中直接保存数据
func (d *KlineDomain) SaveBatch(ctx context.Context, data [][]string, symbol, period string) error {
	mk := make([]*model.Kline, len(data))
	for i, v := range data {
		// 数据封装应该放在model层
		mk[i] = model.NewKline(v, period)
	}

	timestamp := mk[len(mk)-1].Time
	err := d.klineRepo.DeleteGtTime(ctx, timestamp, symbol, period)
	if err != nil {
		return err
	}

	err = d.klineRepo.SaveBatch(ctx, mk, symbol, period)
	if err != nil {
		return err
	}
	logx.Info("KlineDomain | SaveBatch time: ", time.UnixMilli(timestamp).Format("2006-01-02 15:04:05"))
	return nil
}
