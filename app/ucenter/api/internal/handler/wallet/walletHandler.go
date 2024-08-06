package wallet

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic/wallet"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/result"
)

// WalletHandler 获取用户钱包信息
func WalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletReq
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetWalletInfo(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}

// FindWallet 返回用户钱包信息
func FindWallet(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.FindWallet()
		result.HttpResult(r.Context(), w, resp, err)
	}
}

func ResetWalletAddress(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.ResetWalletAddress(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}

func GetAllTransactions(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionReq
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetAllTransactions(&req)
		result.HttpResult(r.Context(), w, resp, err)
	}
}
