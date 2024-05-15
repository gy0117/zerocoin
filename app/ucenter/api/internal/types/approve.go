package types

type ApproveReq struct {
}

type ApproveResp struct {
	Username             string `json:"username"`
	Id                   int64  `json:"id"`
	CreateTime           string `json:"createTime"`
	RealVerified         string `json:"realVerified"`  //是否实名认证
	EmailVerified        string `json:"emailVerified"` //是否有邮箱
	PhoneVerified        string `json:"phoneVerified"` //是否有手机号
	LoginVerified        string `json:"loginVerified"`
	FundsVerified        string `json:"fundsVerified"` //是否有交易密码
	RealAuditing         string `json:"realAuditing"`  // 实名通过的状态
	MobilePhone          string `json:"mobilePhone"`
	Email                string `json:"email"`
	RealName             string `json:"realName"`
	RealNameRejectReason string `json:"realNameRejectReason"`
	IdCard               string `json:"idCard"`
	Avatar               string `json:"avatar"`
	AccountVerified      string `json:"accountVerified"`
}
