package zeroerr

var msgs map[uint32]string

const defaultMsg = "服务器开小差啦,稍后再来试一试"

func init() {
	msgs = make(map[uint32]string)
	msgs[OK] = "SUCCESS"
	msgs[SERVER_COMMON_ERROR] = "服务器开小差啦,稍后再来试一试"
	msgs[REQUEST_PARAM_ERROR] = "参数错误"
	msgs[TOKEN_EXPIRE_ERROR] = "token失效,请重新登陆"
	msgs[DB_ERROR] = "数据库繁忙,请稍后再试"
}

func ParseErrMsg(errCode uint32) string {
	if msg, ok := msgs[errCode]; ok {
		return msg
	} else {
		return defaultMsg
	}
}

func IsErrCode(errCode uint32) bool {
	_, ok := msgs[errCode]
	return ok
}
