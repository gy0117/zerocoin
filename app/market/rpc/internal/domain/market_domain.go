package domain

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"grpc-common/market/types/market"
	"market-rpc/internal/dao"
	"market-rpc/internal/model"
	"market-rpc/internal/repo"
	"math"
	"time"
	"zero-common/operate"
	"zero-common/tools"
)

type MarketDomain struct {
	klineRepo repo.KlineRepo
}

func NewMarketDomain(db *mongo.Database) *MarketDomain {
	return &MarketDomain{
		klineRepo: dao.NewKlineDao(db),
	}
}

// FindSymbolThumbTrend 一个CoinThumb对象就是一个coin的行情数据
func (d *MarketDomain) FindSymbolThumbTrend(coins []*model.ExchangeCoin) []*market.CoinThumb {
	// 从mongodb中找出24h内的数据
	coinThumbs := make([]*market.CoinThumb, len(coins))

	for i, coin := range coins {
		from := tools.ZeroTime()
		end := time.Now().UnixMilli()
		klines, err := d.klineRepo.FindBySymbolTime(context.Background(), coin.Symbol, "1H", from, end, "desc")
		if err != nil {
			coinThumbs[i] = model.DefaultCoinThumb(coin.Symbol)
			logx.Error(err)
			continue
		}
		if len(klines) == 0 {
			coinThumbs[i] = model.DefaultCoinThumb(coin.Symbol)
			continue
		}

		kLen := len(klines)
		trend := make([]float64, kLen)

		var maxVal float64 = math.MinInt64
		var minVal float64 = math.MaxInt64
		var volumes float64 = 0
		var turnovers float64 = 0

		// 计算trend、high、low、volumes、turnovers
		for i := kLen - 1; i >= 0; i-- {
			trend[i] = klines[i].ClosePrice
			high := klines[i].HighestPrice
			if high > maxVal {
				maxVal = high
			}

			low := klines[i].LowestPrice
			if low < minVal {
				minVal = low
			}
			volumes = operate.AddN(volumes, klines[i].Volume, 10)
			turnovers = operate.AddN(turnovers, klines[i].Turnover, 10)
		}

		newKline := klines[0]
		oldKline := klines[kLen-1]
		thumb := newKline.Kline2CoinThumb(coin.Symbol, oldKline)
		thumb.Trend = trend
		thumb.High = maxVal
		thumb.Low = minVal
		thumb.Volume = volumes
		thumb.Turnover = turnovers
		coinThumbs[i] = thumb
	}
	return coinThumbs
}

func (d *MarketDomain) GetHistoryKline(ctx context.Context, symbol string, from int64, to int64, resolution string, period string) ([]*market.History, error) {
	klines, err := d.klineRepo.FindBySymbolTime(ctx, symbol, period, from, to, "asc")
	if err != nil {
		return nil, err
	}
	list := make([]*market.History, len(klines))
	for i, v := range klines {
		h := &market.History{
			Time:   v.Time,
			Open:   v.OpenPrice,
			High:   v.HighestPrice,
			Low:    v.LowestPrice,
			Volume: v.Volume,
			Close:  v.ClosePrice,
		}
		list[i] = h
	}
	return list, nil
}
