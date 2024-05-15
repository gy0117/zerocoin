package zeroerr

import "fmt"

var _ error = (*CodeError)(nil)

type CodeError struct {
	errCode uint32
	errMsg  string
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: ParseErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: SERVER_COMMON_ERROR, errMsg: errMsg}
}

func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}
