package types

type WithdrawReq struct {
	Unit       string  `json:"unit,optional" form:"unit,optional"`             // 币种
	Address    string  `json:"address,optional" form:"address,optional"`       // 提现地址
	Amount     float64 `json:"amount,optional" form:"amount,optional"`         // 提现数量
	Fee        float64 `json:"fee,optional" form:"fee,optional"`               // 手续费
	JyPassword string  `json:"jyPassword,optional" form:"jyPassword,optional"` // 交易密码
	Code       string  `json:"code,optional" form:"code,optional"`             // 验证码
	Page       int64   `json:"page,optional" form:"page,optional"`             // 当前页
	PageSize   int64   `json:"pageSize,optional" form:"pageSize,optional"`     // 每页请求数量
}

type WithdrawWalletInfo struct {
	Unit            string          `json:"unit"`
	Threshold       float64         `json:"threshold"` //阈值
	MinAmount       float64         `json:"minAmount"` //最小提币数量
	MaxAmount       float64         `json:"maxAmount"` //最大提币数量
	MinTxFee        float64         `json:"minTxFee"`  //最小交易手续费
	MaxTxFee        float64         `json:"maxTxFee"`
	NameCn          string          `json:"nameCn"`
	Name            string          `json:"name"`
	Balance         float64         `json:"balance"`
	CanAutoWithdraw string          `json:"canAutoWithdraw"` //true false
	WithdrawScale   int             `json:"withdrawScale"`
	AccountType     int             `json:"accountType"`
	Addresses       []AddressSimple `json:"addresses"`
}

// AddressSimple member_address表中的
type AddressSimple struct {
	Remark  string `json:"remark"`
	Address string `json:"address"`
}
