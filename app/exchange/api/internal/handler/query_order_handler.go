package handler

import (
	"common/result"
	"common/tools"
	"exchange-api/internal/logic"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// QueryHistoryOrders handler --- logic -- domain --- repo(dao)
// 历史委托订单
func QueryHistoryOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := logic.NewQueryOrderLogic(r.Context(), svcCtx)
		resp, err := l.QueryHistoryOrders(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// QueryCurrentOrders 当前委托订单
func QueryCurrentOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := logic.NewQueryOrderLogic(r.Context(), svcCtx)
		resp, err := l.QueryCurrentOrders(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

func QueryCompleteOrders(serviceCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := logic.NewQueryOrderLogic(r.Context(), serviceCtx)
		resp, err := l.QueryCompleteOrders(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
