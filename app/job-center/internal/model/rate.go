package model

type ExchangeRate struct {
	UsdCny string `json:"usdCny"`
}

type OkxRateResp struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []*ExchangeRate `json:"data"`
}
