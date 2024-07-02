package market

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"grpc-common/market/types/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketLogic) SymbolThumbTrend(req *types.MarketRequest) (resp []*types.CoinThumbResponse, err error) {
	var thumbs []*market.CoinThumb
	// 有缓存，先从缓存中拿，没有缓存，请求rpc
	thumb := l.svcCtx.Processor.GetThumb()
	isCache := false
	if thumb != nil {
		switch thumb.(type) {
		case []*market.CoinThumb:
			thumbs = thumb.([]*market.CoinThumb)
			isCache = true
		}
	}

	if !isCache {
		ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
		defer cancel()
		// api和rpc模块中的Req参数是类似的
		result, err := l.svcCtx.MarketRpc.FindSymbolThumbTrend(ctx, &market.MarketRequest{
			Ip:         req.Ip,
			Symbol:     req.Symbol,
			Unit:       req.Unit,
			From:       req.From,
			To:         req.To,
			Resolution: req.Resolution,
		})
		if err != nil {
			return nil, err
		}
		thumbs = result.List
	}

	// rpc数据转换为api数据
	resp = make([]*types.CoinThumbResponse, len(thumbs))
	for i, v := range thumbs {
		resp[i] = &types.CoinThumbResponse{
			Symbol:       v.Symbol,
			OpenPrice:    v.Open,
			HighPrice:    v.High,
			LowPrice:     v.Low,
			ClosePrice:   v.Close,
			Chg:          v.Chg,
			Change:       v.Change,
			Volume:       v.Volume,
			Turnover:     v.Turnover,
			LastDayClose: v.LastDayClose,
			UsdRate:      v.UsdRate,
			BaseUsdRate:  v.BaseUsdRate,
			Zone:         int(v.Zone),
			Trend:        v.Trend,
		}
	}
	return
}

func (l *MarketLogic) SymbolThumb(req *types.MarketRequest) (resp []*types.CoinThumbResponse, err error) {
	var thumbs []*market.CoinThumb
	thumb := l.svcCtx.Processor.GetThumb()
	if thumb != nil {
		switch thumb.(type) {
		case []*market.CoinThumb:
			thumbs = thumb.([]*market.CoinThumb)
		}
	}

	resp = make([]*types.CoinThumbResponse, len(thumbs))
	for i, v := range thumbs {
		resp[i] = &types.CoinThumbResponse{
			Symbol:       v.Symbol,
			OpenPrice:    v.Open,
			HighPrice:    v.High,
			LowPrice:     v.Low,
			ClosePrice:   v.Close,
			Chg:          v.Chg,
			Change:       v.Change,
			Volume:       v.Volume,
			Turnover:     v.Turnover,
			LastDayClose: v.LastDayClose,
			UsdRate:      v.UsdRate,
			BaseUsdRate:  v.BaseUsdRate,
			Zone:         int(v.Zone),
			Trend:        v.Trend,
		}
	}

	return
}

func (l *MarketLogic) SymbolInfo(req *types.MarketRequest) (*types.ExchangeCoinResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// 调用rpc
	marketReq := &market.MarketRequest{
		Ip:     req.Ip,
		Symbol: req.Symbol,
	}

	result, err := l.svcCtx.MarketRpc.FindSymbolInfo(ctx, marketReq)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SymbolInfo---result: %+v\n", result)
	resp := &types.ExchangeCoinResp{}
	if err = copier.Copy(resp, result); err != nil {
		logx.Error(err)
		return nil, err
	}
	return resp, nil
}

func (l *MarketLogic) CoinInfo(req *types.MarketRequest) (*types.Coin, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()

	// 调用rpc
	in := &market.MarketRequest{
		Ip:   req.Ip,
		Unit: req.Unit,
	}
	coin, err := l.svcCtx.MarketRpc.FindCoinInfo(ctx, in)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	resp := &types.Coin{}
	if err = copier.Copy(resp, coin); err != nil {
		logx.Error(err)
		return nil, err
	}
	return resp, nil
}

func (l *MarketLogic) GetHistoryKline(req *types.MarketRequest) (*types.HistoryKline, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()

	// 调用rpc
	in := &market.MarketRequest{
		Ip:         req.Ip,
		Symbol:     req.Symbol,
		From:       req.From,
		To:         req.To,
		Resolution: req.Resolution,
	}
	kline, err := l.svcCtx.MarketRpc.GetHistoryKline(ctx, in)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	resp := make([][]any, len(kline.List))
	for i, v := range kline.List {
		content := make([]any, 6)
		content[0] = v.GetTime()
		content[1] = v.GetOpen()
		content[2] = v.GetHigh()
		content[3] = v.GetLow()
		content[4] = v.GetClose()
		content[5] = v.GetVolume()

		resp[i] = content
	}

	return &types.HistoryKline{
		List: resp,
	}, nil
}
