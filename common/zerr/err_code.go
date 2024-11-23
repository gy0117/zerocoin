package zerr

const OK uint32 = 200

// 前三位表示业务，后三位表示功能
// 全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REQUEST_PARAM_ERROR uint32 = 100002
const TOKEN_EXPIRE_ERROR uint32 = 100003
const TOKEN_GENERATE_ERROR uint32 = 100004
const DB_ERROR uint32 = 100005

// ucenter模块 200xxx
const USER_REGISTER_ERROR = 200001       // 用户注册失败
const USER_HAS_REGISTERED_ERROR = 200002 // 用户已经注册
const USER_VERIFY_CODE_ERROR = 200003    // 验证码填写错误

const USER_LOGIN_ERROR = 200004           // 登录失败
const USER_PHONE_NOT_EXIST_ERROR = 200005 // 手机号不存在
const USER_PASSWORD_ERROR = 200006        // 密码错误

const FIND_USER_ERROR = 200007 // 查找用户失败

const FIND_WALLET_ERROR = 200008          // 查找钱包失败
const RESET_WALLET_ADDRESS_ERROR = 200009 // 重置钱包地址失败

const GET_TRANSACTIONS_ERROR = 200010 // 获取交易失败
const GET_ADDRESS_ERROR = 200011      // 获取地址失败
const WITHDRAW_ERROR = 200012         // 提现失败
const WITHDRAW_FIND_RECORD = 200013   // 查询提现记录失败

// market模块 300xxx
const MARKET_FIND_SYMBOL_ERROR = 300001   // 查询symbol失败
const MARKET_FIND_COIN_ERROR = 300002     // 查询coin失败
const MARKET_HISTORY_KLINE_ERROR = 300003 // 查询历史kline失败

// exchange模块 400xxx
const EXCHANGE_GET_HISTORY_ORDER_ERROR = 400001 // 获取历史订单失败
const EXCHANGE_GET_CURRENT_ORDER_ERROR = 400002 // 获取当前订单失败
const EXCHANGE_ADD_ORDER_ERROR = 400003         // 添加订单失败
const EXCHANGE_FIND_ORDER_ERROR = 400004        // 查询订单失败
const EXCHANGE_CANCEL_ORDER_ERROR = 400005      // 取消订单

// job-center模块 500xxx
