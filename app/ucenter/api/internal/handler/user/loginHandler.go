package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic/user"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/result"
	"zero-common/tools"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		// 获取ip
		req.Ip = tools.GetRemoteClientIp(r)

		env := r.Header.Get("env")
		req.Env = env

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
