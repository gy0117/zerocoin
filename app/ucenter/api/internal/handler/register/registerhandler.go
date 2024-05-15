package register

import (
	"net/http"
	"ucenter-api/internal/logic/register"
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
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取ip
		req.Ip = tools.GetRemoteClientIp(r)

		// 获取env
		env := r.Header.Get("env")
		req.Env = env

		l := register.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}
