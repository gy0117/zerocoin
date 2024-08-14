package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/result"
)

func SendCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CodeReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.SendCode(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
