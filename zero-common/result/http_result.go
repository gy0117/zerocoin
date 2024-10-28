package result

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
	"zero-common/zerr"
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

func HttpResult2(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, Success(resp))
	} else {
		errCode := zerr.SERVER_COMMON_ERROR
		errMsg := zerr.GetDefaultMsg()

		// 判断错误是自定义的，还是grpc的
		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zerr.CodeError); ok { // api err 自定义错误类型
			errCode = e.GetErrCode()
			errMsg = e.GetErrMsg()
		} else {
			if grpcStatus, ok := status.FromError(causeErr); ok { // rpc err 自定义错误类型 或者 grpc err
				grpcCode := uint32(grpcStatus.Code())
				if zerr.IsCodeErr(grpcCode) { // 自定义错误返回出去，grpc自己的错误就不要返回真实信息给前端了
					errCode = grpcCode
					errMsg = grpcStatus.Message()
				}
			}
		}
		logx.WithContext(r.Context()).Errorf("API-ERR | err: %+v", err)
		httpx.WriteJson(w, http.StatusBadRequest, Error(errCode, errMsg))
	}
}

// 参数错误返回
func ParamErrorResult(w http.ResponseWriter, r *http.Request, err error) {
	errMsg := fmt.Sprintf("%s ,%s", zerr.ParseErrMsg(zerr.REQUEST_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(zerr.REQUEST_PARAM_ERROR, errMsg))
}
