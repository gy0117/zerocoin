package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-rpc/internal/domain"
	"market-rpc/internal/svc"
	"time"
	"zero-common/zerr"
)

var ErrMarketFindSymbol = zerr.NewCodeErr(zerr.MARKET_FIND_SYMBOL_ERROR)
var ErrMarketFindCoin = zerr.NewCodeErr(zerr.MARKET_FIND_COIN_ERROR)
var ErrMarketHistoryKline = zerr.NewCodeErr(zerr.MARKET_FIND_COIN_ERROR)

type MarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeCoinDomain *domain.ExchangeCoinDomain
	marketDomain       *domain.MarketDomain
	coinDomain         *domain.CoinDomain
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		exchangeCoinDomain: domain.NewExchangeCoinDomain(svcCtx.DB),
		marketDomain:       domain.NewMarketDomain(svcCtx.MongoClient.Db),
		coinDomain:         domain.NewCoinDomain(svcCtx.DB),
	}
}

func (l *MarketLogic) FindSymbolThumbTrend(in *market.MarketRequest) (*market.CoinThumbResponse, error) {
	// 1. 先找到支持哪些coin
	coins := l.exchangeCoinDomain.FindCoinVisible(l.ctx)

	// 2. 找到所支持的coin的行情
	coinThumbs := l.marketDomain.FindSymbolThumbTrend(coins)

	resp := &market.CoinThumbResponse{
		List: coinThumbs,
	}
	return resp, nil
}

func (l *MarketLogic) FindSymbolInfo(in *market.MarketRequest) (*market.ExchangeCoin, error) {
	// 1. 查询exchange_coin表中的数据，根据symbol
	symbol := in.GetSymbol()
	coin, err := l.exchangeCoinDomain.FindBySymbol(l.ctx, symbol)
	if err != nil {
		return nil, errors.Wrapf(ErrMarketFindSymbol, "market find symbol: %s", symbol)
	}

	resp := &market.ExchangeCoin{
		Id:               coin.Id,
		Symbol:           coin.Symbol,
		BaseCoinScale:    coin.BaseCoinScale,
		BaseSymbol:       coin.BaseSymbol,
		CoinScale:        coin.CoinScale,
		Enable:           coin.Enable,
		Fee:              coin.Fee,
		Sort:             coin.Sort,
		EnableMarketBuy:  coin.EnableMarketBuy,
		EnableMarketSell: coin.EnableMarketSell,
		MinSellPrice:     coin.MinSellPrice,
		Flag:             coin.Flag,
		MaxTradingOrder:  coin.MaxTradingOrder,
		MaxTradingTime:   coin.MaxTradingTime,
		MinTurnover:      int64(coin.MinTurnover),
		ClearTime:        coin.ClearTime,
		EndTime:          coin.EndTime,
		Exchangeable:     coin.Exchangeable,
		MaxBuyPrice:      coin.MaxBuyPrice,
		MaxVolume:        coin.MaxVolume,
		MinVolume:        coin.MinVolume,
		PublishAmount:    coin.PublishAmount,
		PublishPrice:     coin.PublishPrice,
		PublishType:      coin.PublishType,
		RobotType:        coin.RobotType,
		StartTime:        coin.StartTime,
		Visible:          coin.Visible,
		Zone:             coin.Zone,
		CoinSymbol:       coin.CoinSymbol,
	}
	return resp, nil
}

func (l *MarketLogic) FindCoinInfo(in *market.MarketRequest) (*market.Coin, error) {
	coin, err := l.coinDomain.FindCoinInfo(l.ctx, in.Unit)

	if err != nil {
		return nil, errors.Wrapf(ErrMarketFindCoin, "market find coin: %s", in.Unit)
	}
	resp := &market.Coin{}
	_ = copier.Copy(resp, coin)
	return resp, nil
}

func (l *MarketLogic) GetHistoryKline(in *market.MarketRequest) (*market.HistoryResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()
	period := "1H"
	if in.GetResolution() == "1" {
		period = "1m"
	} else if in.GetResolution() == "5" {
		period = "5m"
	} else if in.GetResolution() == "15" {
		period = "15m"
	} else if in.GetResolution() == "30" {
		period = "30m"
	} else if in.GetResolution() == "60" {
		period = "1H"
	} else if in.GetResolution() == "1D" {
		period = "1D"
	} else if in.GetResolution() == "1W" {
		period = "1W"
	} else if in.GetResolution() == "1M" {
		period = "1M"
	}

	histories, err := l.marketDomain.GetHistoryKline(ctx, in.GetSymbol(), in.GetFrom(), in.GetTo(), in.GetResolution(), period)
	if err != nil {
		return nil, errors.Wrapf(ErrMarketHistoryKline, "market getHistoryKline, symbol: %s, from: %d, to: %d, period: %s", in.GetSymbol(), in.GetFrom(), in.GetTo(), period)
	}

	return &market.HistoryResp{
		List: histories,
	}, nil
}

func (l *MarketLogic) FindExchangeCoinVisible(in *market.MarketRequest) (*market.ExchangeCoinResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()
	coins := l.exchangeCoinDomain.FindCoinVisible(ctx)

	var list []*market.ExchangeCoin

	_ = copier.Copy(&list, coins)
	return &market.ExchangeCoinResp{
		List: list,
	}, nil
}

func (l *MarketLogic) FindCoinByCoinId(in *market.MarketRequest) (*market.Coin, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()

	coin, err := l.coinDomain.FindCoinByCoinId(ctx, in.CoinId)
	if err != nil {
		return nil, errors.Wrapf(ErrMarketFindCoin, "market findCoinByCoinId, coindId: %d", in.CoinId)
	}
	resp := &market.Coin{}
	_ = copier.Copy(resp, coin)
	return resp, nil
}
