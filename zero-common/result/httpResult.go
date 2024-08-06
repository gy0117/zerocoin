package result

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type ZeroCode int64

const ErrorCode ZeroCode = -1
const SuccessCode ZeroCode = 0

const SuccessMsg string = "success"

type Result struct {
	Code    ZeroCode `json:"code"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) Success(data any) {
	r.Code = SuccessCode
	r.Message = SuccessMsg
	r.Data = data
}

func (r *Result) Fail(code ZeroCode, message string) {
	r.Code = code
	r.Message = message
}

func (r *Result) Deal(data any, err error) {
	if err != nil {
		r.Fail(-1, err.Error())
	}
}

// HttpResult 返回给前端
func HttpResult(ctx context.Context, w http.ResponseWriter, resp interface{}, err error) {
	result := NewResult()
	code := SuccessCode
	msg := SuccessMsg
	if err != nil {
		code = ErrorCode
		msg = err.Error()
	} else {
		if resp != nil {
			result.Data = resp
		}
	}
	result.Code = code
	result.Message = msg

	httpx.OkJsonCtx(ctx, w, result)
}
