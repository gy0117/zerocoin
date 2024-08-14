package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"market-api/internal/logic"
	"market-api/internal/svc"
	"market-api/internal/types"
	"net/http"
	"zero-common/result"
	"zero-common/tools"
)

func ExchangeRateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RateRequest
		if err := httpx.ParsePath(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		req.Ip = tools.GetRemoteClientIp(r)

		l := logic.NewExchangeRateLogic(r.Context(), svcCtx)
		resp, err := l.UsdRate(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
