package middleware

import (
	common "common/result"
	"common/tools"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func Auth(secret string, f func(string) bool) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			result := common.NewResult()

			token := r.Header.Get("X-Auth-Token")
			if token == "" {
				result.Fail(4000, "not login")
				httpx.WriteJson(w, 200, result)
				return
			}
			userId, err := tools.ParseToken(token, secret)
			if err != nil {
				result.Fail(4000, "token is error")
				httpx.WriteJson(w, 200, result)
				return
			}

			// check black list
			newToken := fmt.Sprintf("auth:token:%s", token)
			if f(newToken) {
				result.Fail(4000, "token is error")
				httpx.WriteJson(w, 200, result)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			r = r.WithContext(ctx)
			next(w, r)
		}
	}
}
