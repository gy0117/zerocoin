package result

type ResponseSuccessData struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *ResponseSuccessData {
	return &ResponseSuccessData{
		Code: 200,
		Msg:  "OK",
		Data: data,
	}
}

type ResponseErrorData struct {
	Code uint32 `json:"code"`
	Msg  string `json:"message"`
}

func Error(errCode uint32, errMsg string) *ResponseErrorData {
	return &ResponseErrorData{
		Code: errCode,
		Msg:  errMsg,
	}
}
