package types

type WalletReq struct {
	CoinName string `json:"coinName,optional" path:"coinName,optional"`
	Unit     string `json:"unit,optional" form:"unit,optional"`
}

type UserWallet struct {
	Id             int64   `json:"id" from:"id"`
	Address        string  `json:"address" from:"address"`
	Balance        float64 `json:"balance" from:"balance"`
	FrozenBalance  float64 `json:"frozenBalance" from:"frozenBalance"`
	ReleaseBalance float64 `json:"releaseBalance" from:"releaseBalance"`
	IsLock         int     `json:"isLock" from:"isLock"`
	UserId         int64   `json:"userId" from:"userId"`
	Version        int     `json:"version" from:"version"`
	Coin           Coin    `json:"coin" from:"coinId"`
	ToReleased     float64 `json:"toReleased" from:"toReleased"`
}

type Coin struct {
	Id                int     `json:"id" from:"id"`
	Name              string  `json:"name" from:"name"`
	CanAutoWithdraw   int     `json:"canAutoWithdraw" from:"canAutoWithdraw"`
	CanRecharge       int     `json:"canRecharge" from:"canRecharge"`
	CanTransfer       int     `json:"canTransfer" from:"canTransfer"`
	CanWithdraw       int     `json:"canWithdraw" from:"canWithdraw"`
	CnyRate           float64 `json:"cnyRate" from:"cnyRate"`
	EnableRpc         int     `json:"enableRpc" from:"enableRpc"`
	IsPlatformCoin    int     `json:"isPlatformCoin" from:"isPlatformCoin"`
	MaxTxFee          float64 `json:"maxTxFee" from:"maxTxFee"`
	MaxWithdrawAmount float64 `json:"maxWithdrawAmount" from:"maxWithdrawAmount"`
	MinTxFee          float64 `json:"minTxFee" from:"minTxFee"`
	MinWithdrawAmount float64 `json:"minWithdrawAmount" from:"minWithdrawAmount"`
	NameCn            string  `json:"nameCn" from:"nameCn"`
	Sort              int     `json:"sort" from:"sort"`
	Status            int     `json:"status" from:"status"`
	Unit              string  `json:"unit" from:"unit"`
	UsdRate           float64 `json:"usdRate" from:"usdRate"`
	WithdrawThreshold float64 `json:"withdrawThreshold" from:"withdrawThreshold"`
	HasLegal          int     `json:"hasLegal" from:"hasLegal"`
	ColdWalletAddress string  `json:"coldWalletAddress" from:"coldWalletAddress"`
	MinerFee          float64 `json:"minerFee" from:"minerFee"`
	WithdrawScale     int     `json:"withdrawScale" from:"withdrawScale"`
	AccountType       int     `json:"accountType" from:"accountType"`
	DepositAddress    string  `json:"depositAddress" from:"depositAddress"`
	Infolink          string  `json:"infolink" from:"infolink"`
	Information       string  `json:"information" from:"information"`
	MinRechargeAmount float64 `json:"minRechargeAmount" from:"minRechargeAmount"`
}

type TransactionReq struct {
	PageNo    int64  `json:"pageNo" form:"pageNo"`                         // 当前页
	PageSize  int64  `json:"pageSize" form:"pageSize"`                     // 每页显示数量
	StartTime string `json:"startTime,optional" form:"startTime,optional"` // 开始时间
	EndTime   string `json:"endTime,optional" form:"endTime,optional"`     // 结束时间
	Symbol    string `json:"symbol,optional" form:"symbol,optional"`       // 币种名称 比如 BTC
	Type      string `json:"type,optional" form:"type,optional"`           // 0 RECHARGE 充值类型 1 WITHDRAW 提现类型
}

type UserTransaction struct {
	Id          int64   `json:"id" from:"id"`
	Address     string  `json:"address" from:"address"`
	Amount      float64 `json:"amount" from:"amount"`
	CreateTime  string  `json:"createTime" from:"createTime"`
	Fee         float64 `json:"fee" from:"fee"`
	Flag        int     `json:"flag" from:"flag"`
	UserId      int64   `json:"userId" from:"userId"`
	Symbol      string  `json:"symbol" from:"symbol"`
	Type        string  `json:"type" from:"type"`
	DiscountFee string  `json:"discountFee" from:"discountFee"`
	RealFee     string  `json:"realFee" from:"realFee"`
}
