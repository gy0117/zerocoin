// Code generated by goctl. DO NOT EDIT.
package types

type Request struct {
	Username     string      `json:"username, optional"`
	Password     string      `json:"password,optional"`
	Captcha      *CaptchaReq `json:"captcha,optional"`
	Phone        string      `json:"phone,optional"`
	Promotion    string      `json:"promotion,optional"`
	Code         string      `json:"code,optional"`
	Country      string      `json:"country,optional"`
	SuperPartner string      `json:"superPartner,optional"`
	Ip           string      `json:"ip,optional"`
	Env          string      `json:"env,optional"`
}

type CaptchaReq struct {
	Server string `json:"server"`
	Token  string `json:"token"`
}

type Response struct {
	Message string `json:"message"`
}

type CodeReq struct {
	Country string `json:"country"`
	Phone   string `json:"phone"`
}

type CodeResp struct {
	SmsCode string `json:"smsCode"`
}

// 注册 resp
type RegisterResp struct {
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	RefreshAfter int64  `json:"refreshAfter"`
}