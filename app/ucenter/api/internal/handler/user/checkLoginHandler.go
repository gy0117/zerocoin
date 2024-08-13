package user

import (
	"net/http"
	"ucenter-api/internal/logic/user"
	"ucenter-api/internal/svc"
	"zero-common/result"
)

func CheckLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		l := user.NewLoginLogic(r.Context(), svcCtx)
		isOk, err := l.CheckLogin(token)
		result.HttpResult2(w, r, isOk, err)
	}
}
