package login

import (
	"net/http"
	"ucenter-api/internal/logic/login"
	"ucenter-api/internal/svc"
	"zero-common/result"
)

func CheckLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		l := login.NewLoginLogic(r.Context(), svcCtx)
		isOk, err := l.CheckLogin(token)
		result.HttpResult(r.Context(), w, isOk, err)
	}
}
