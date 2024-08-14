package handler

import (
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"zero-common/result"
)

func SecuritySetting(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewApproveLogic(r.Context(), svcCtx)
		resp, err := l.CheckSecuritySetting()
		result.HttpResult2(w, r, resp, err)
	}
}
