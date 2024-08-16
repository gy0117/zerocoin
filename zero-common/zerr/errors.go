package zerr

import (
	"fmt"
)

var _ error = (*CodeError)(nil)

// 错误类型定义
type CodeError struct {
	errCode uint32
	errMsg  string
}

// 返回给前端展示的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// 返回给前端展示的错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", e.errCode, e.errMsg)
}

func NewCodeMsgErr(errCode uint32, errMsg string) *CodeError {
	return &CodeError{
		errCode: errCode,
		errMsg:  errMsg,
	}
}

func NewCodeErr(errCode uint32) *CodeError {
	return &CodeError{
		errCode: errCode,
		errMsg:  ParseErrMsg(errCode),
	}
}

func NewMsgErr(errMsg string) *CodeError {
	return &CodeError{
		errCode: SERVER_COMMON_ERROR,
		errMsg:  errMsg,
	}
}
