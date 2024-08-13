package wallet

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic/wallet"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/result"
)

// 获取用户钱包信息
func GetWalletWithCoin(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletReq
		if err := httpx.ParsePath(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetWalletInfo(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// 返回用户钱包信息
func GetWallet(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.FindWallet()
		result.HttpResult2(w, r, resp, err)
	}
}

// 重置钱包地址
func ResetWalletAddress(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}

		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.ResetWalletAddress(&req)
		result.HttpResult2(w, r, resp, err)
	}
}

// 获取所有交易
func GetAllTransactions(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransactionReq
		if err := httpx.ParseForm(r, &req); err != nil {
			result.ParamErrorResult(w, r, err)
			return
		}
		l := wallet.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetAllTransactions(&req)
		result.HttpResult2(w, r, resp, err)
	}
}
