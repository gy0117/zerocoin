package logic

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"job-center/internal/config"
	"job-center/internal/model"
	"time"
	"zero-common/tools"
)

// 1 USDT 对应多少CNY
const keyUsdt2Cny = "USDT::CNY::RATE"

type Rate struct {
	conf  config.OkxConfig
	cache cache.Cache
}

func (r *Rate) Do() {
	logx.Info("获取法币汇率")

	go r.getCny2Usd()
}

// 获取人民币对美元汇率
func (r *Rate) getCny2Usd() {
	url := r.conf.Host + "/api/v5/market/exchange-rate"

	timestamp := tools.ISO(time.Now())
	sign := tools.ComputeHmacSha256(timestamp+"GET"+"/api/v5/market/exchange-rate", r.conf.SecretKey)
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = r.conf.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = r.conf.Passphrase

	// 访问外网需要添加代理
	bytes, err := tools.GetWithHeader(url, header, r.conf.Proxy)
	if err != nil {
		logx.Error(err)
		return
	}

	resp := &model.OkxRateResp{}
	if err := json.Unmarshal(bytes, resp); err != nil {
		logx.Error(err)
		return
	}
	// 存到redis
	if len(resp.Data) > 0 {
		r.saveCache(resp.Data[0])
	}
}

func (r *Rate) saveCache(rate *model.ExchangeRate) {
	err := r.cache.Set(keyUsdt2Cny, rate.UsdCny)
	if err != nil {
		logx.Error(err)
	}
}

func NewRate(conf config.OkxConfig, redis cache.Cache) *Rate {
	return &Rate{
		conf:  conf,
		cache: redis,
	}
}
