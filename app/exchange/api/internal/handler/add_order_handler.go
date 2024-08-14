package handler

import (
	"exchange-api/internal/logic"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"zero-common/result"
	"zero-common/tools"
)

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
		l := logic.NewCreateOrderLogic(r.Context(), svcCtx)
		orderId, err := l.AddOrder(&req)
		result.HttpResult2(w, r, orderId, err)
	}
}
