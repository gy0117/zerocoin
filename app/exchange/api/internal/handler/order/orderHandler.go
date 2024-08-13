package order

import (
	"errors"
	"exchange-api/internal/logic/order"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"zero-common/result"
	"zero-common/tools"
)

// GetHistoryOrders handler --- logic -- domain --- repo(dao)
// 历史委托订单
func GetHistoryOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetHistoryOrders(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// GetCurrentOrders 当前委托订单
func GetCurrentOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentOrders(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// AddOrder 发布委托
func AddOrder(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		if !req.IsValid() {
			result.ParamErrorResult(w, r, errors.New("订单参数错误"))
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		orderId, err := l.AddOrder(&req)
		result.HttpResult2(w, r, orderId, err)
	}
}
