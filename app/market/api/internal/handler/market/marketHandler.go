package market

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"market-api/internal/logic/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"net/http"
	"zero-common/result"
	"zero-common/tools"
)

// SymbolThumbTrendHandler 获取币种行情
func SymbolThumbTrendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarketRequest
		req.Ip = tools.GetRemoteClientIp(r)

		l := market.NewMarketLogic(r.Context(), svcCtx)
		resp, err := l.SymbolThumbTrend(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// SymbolThumbHandler 币币交易部分
func SymbolThumbHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarketRequest
		req.Ip = tools.GetRemoteClientIp(r)

		l := market.NewMarketLogic(r.Context(), svcCtx)
		resp, err := l.SymbolThumb(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// SymbolInfoHandler 交易币详情
func SymbolInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarketRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		req.Ip = tools.GetRemoteClientIp(r)

		l := market.NewMarketLogic(r.Context(), svcCtx)
		resp, err := l.SymbolInfo(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// CoinInfoHandler 货币详情
func CoinInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarketRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)

		l := market.NewMarketLogic(r.Context(), svcCtx)
		resp, err := l.CoinInfo(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// HistoryHandler 币币交易-k线图
func HistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarketRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)

		l := market.NewMarketLogic(r.Context(), svcCtx)
		resp, err := l.GetHistoryKline(&req)
		result.HttpResult2(w, r, resp.List, err) // 返回[][]any
	}
}
