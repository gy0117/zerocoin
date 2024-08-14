// Code generated by goctl. DO NOT EDIT.
// Source: market.proto

package server

import (
	"context"
	"grpc-common/market/types/market"
	"market-rpc/internal/logic"
	"market-rpc/internal/svc"
)

type MarketServer struct {
	svcCtx *svc.ServiceContext
	market.UnimplementedMarketServer
}

func NewMarketServer(svcCtx *svc.ServiceContext) *MarketServer {
	return &MarketServer{
		svcCtx: svcCtx,
	}
}

func (ms *MarketServer) FindSymbolThumbTrend(ctx context.Context, in *market.MarketRequest) (*market.CoinThumbResponse, error) {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.FindSymbolThumbTrend(in)
}

func(ms *MarketServer) FindSymbolInfo(ctx context.Context, in *market.MarketRequest) (*market.ExchangeCoin, error)  {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.FindSymbolInfo(in)
}

// rpc error: code = Unimplemented desc = method FindCoinInfo not implemented
func(ms *MarketServer) FindCoinInfo(ctx context.Context, in *market.MarketRequest) (*market.Coin, error)  {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.FindCoinInfo(in)
}

func(ms *MarketServer) GetHistoryKline(ctx context.Context, in *market.MarketRequest) (*market.HistoryResp, error)  {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.GetHistoryKline(in)
}

func(ms *MarketServer) FindExchangeCoinVisible(ctx context.Context, in *market.MarketRequest) (*market.ExchangeCoinResp, error)  {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.FindExchangeCoinVisible(in)
}

func(ms *MarketServer) FindCoinByCoinId(ctx context.Context, in *market.MarketRequest) (*market.Coin, error)  {
	l := logic.NewMarketLogic(ctx, ms.svcCtx)
	return l.FindCoinByCoinId(in)
}