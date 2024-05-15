package types

type LoginReq struct {
	UserName string      `json:"username"`
	Password string      `json:"password"`
	Captcha  *CaptchaReq `json:"captcha,optional"`
	Ip       string      `json:"ip,optional"`
	Env      string      `json:"env,optional"`
}

type LoginResp struct {
	UserName      string `json:"username"`
	Token         string `json:"token"`
	MemberLevel   string `json:"memberLevel"`
	RealName      string `json:"realName"`
	Country       string `json:"country"`
	Avatar        string `json:"avatar"`
	PromotionCode string `json:"promotionCode"`
	Id            int64  `json:"id"`
	LoginCount    int    `json:"loginCount"`
	SuperPartner  string `json:"superPartner"`
	MemberRate    int    `json:"memberRate"`
}
