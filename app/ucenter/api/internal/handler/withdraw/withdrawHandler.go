package withdraw

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic/withdraw"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/result"
)

func GetSupportedCoinInfo(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := withdraw.NewWithdrawLogic(r.Context(), svcCtx)
		resp, err := l.GetSupportedCoinInfo()
		result.HttpResult(r.Context(), w, resp, err)
	}
}

func SendCode(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := withdraw.NewWithdrawLogic(r.Context(), svcCtx)
		err := l.SendCode()
		result.HttpResult(r.Context(), w, nil, err)
	}
}

// Withdraw 提现
func Withdraw(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := withdraw.NewWithdrawLogic(r.Context(), svcCtx)
		err := l.Withdraw(&req)
		result.HttpResult(r.Context(), w, nil, err)
	}
}

// Record 提现记录
func Record(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := withdraw.NewWithdrawLogic(r.Context(), svcCtx)
		resp, err := l.Record(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}
