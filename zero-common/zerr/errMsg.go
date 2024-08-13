package zerr

const defaultMsg = "服务器开小差啦,稍后再来试一试"

func GetDefaultMsg() string {
	return defaultMsg
}

var msg map[uint32]string

func init() {
	msg = make(map[uint32]string)
	msg[OK] = "SUCCESS"
	msg[SERVER_COMMON_ERROR] = "服务器开小差啦，请稍后再来试一试～"
	msg[REQUEST_PARAM_ERROR] = "请求参数错误哦"
	msg[TOKEN_EXPIRE_ERROR] = "token过期啦，请重新登陆"
	msg[TOKEN_GENERATE_ERROR] = "token生成失败"
	msg[DB_ERROR] = "数据库开小差啦，请稍后再试一试～"

	// 注册
	msg[USER_REGISTER_ERROR] = "注册失败"
	msg[USER_HAS_REGISTERED_ERROR] = "该用户已被注册"
	msg[USER_VERIFY_CODE_ERROR] = "验证码填写错误"

	// 登录
	msg[USER_LOGIN_ERROR] = "登录失败"
	msg[USER_PHONE_NOT_EXIST_ERROR] = "手机号不存在"
	msg[USER_PASSWORD_ERROR] = "密码"

	// user
	msg[FIND_USER_ERROR] = "查找用户失败"
	// 钱包
	msg[FIND_WALLET_ERROR] = "查找钱包失败"
	msg[RESET_WALLET_ADDRESS_ERROR] = "重置钱包地址失败"
	msg[GET_TRANSACTIONS_ERROR] = "获取交易记录失败"
	msg[GET_ADDRESS_ERROR] = "获取地址失败"
	msg[WITHDRAW_ERROR] = "提现失败"
	msg[WITHDRAW_FIND_RECORD] = "查询提现记录失败"

	// market
	msg[MARKET_FIND_SYMBOL_ERROR] = "查询symbol失败"
	msg[MARKET_FIND_COIN_ERROR] = "查询coin失败"
	msg[MARKET_HISTORY_KLINE_ERROR] = "查询历史kline失败"

	// order
	msg[EXCHANGE_GET_HISTORY_ORDER_ERROR] = "获取历史订单失败"
	msg[EXCHANGE_GET_CURRENT_ORDER_ERROR] = "获取当前订单失败"
	msg[EXCHANGE_ADD_ORDER_ERROR] = "添加订单失败"
	msg[EXCHANGE_FIND_ORDER_ERROR] = "查询订单失败"
	msg[EXCHANGE_CANCEL_ORDER_ERROR] = "取消订单失败"
}

// ParseErrMsg 根据错误码，获取错误信息
func ParseErrMsg(ecode uint32) string {
	if v, ok := msg[ecode]; ok {
		return v
	}
	return defaultMsg
}

// IsCodeErr 是否是自定义错误码
func IsCodeErr(ecode uint32) bool {
	_, ok := msg[ecode]
	return ok
}
