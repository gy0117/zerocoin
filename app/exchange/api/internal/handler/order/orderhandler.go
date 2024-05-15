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
func GetHistoryOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetHistoryOrders(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}

func GetCurrentOrders(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentOrders(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}

func AddOrder(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if !req.IsValid() {
			httpx.ErrorCtx(r.Context(), w, errors.New("订单参数错误"))
			return
		}
		req.Ip = tools.GetRemoteClientIp(r)
		l := order.NewOrderLogic(r.Context(), svcCtx)
		orderId, err := l.AddOrder(&req)
		result.HttpResult(r.Context(), w, orderId, err)
	}
}
