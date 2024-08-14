package handler

import (
	"net/http"
	"ucenter-api/internal/logic"
	"zero-common/result"
	"zero-common/tools"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		// 获取ip
		req.Ip = tools.GetRemoteClientIp(r)

		// 获取env
		env := r.Header.Get("env")
		req.Env = env

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
