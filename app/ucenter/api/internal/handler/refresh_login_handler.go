package handler

import (
	"common/result"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
)

func RefreshTokenHandler(serviceCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken := r.Header.Get("X-Refresh-Token")
		l := logic.NewLoginLogic(r.Context(), serviceCtx)

		accessToken, err := l.RefreshToken(refreshToken)
		result.HttpResult2(w, r, accessToken, err)
	}
}
