package domain

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"strings"
	"zero-common/tools"
)

// 1 USDT 对应多少CNY
const keyUsdt2Cny = "USDT::CNY::RATE"

type ExchangeRateDomain struct {
	cache cache.Cache
}

func NewExchangeRateDomain(db cache.Cache) *ExchangeRateDomain {
	return &ExchangeRateDomain{
		cache: db,
	}
}

// UsdRate 从redis中查询，在job-center中写一个定时任务查询实时汇率，存入redis
func (d *ExchangeRateDomain) UsdRate(unit string) float64 {
	unit = strings.ToUpper(unit)
	if "CNY" == unit {
		return d.getCny2USDRate()
	}
	if "JPY" == unit {
		return 137
	}
	return 0
}

// 获取 1USD 2 cny的汇率
func (d *ExchangeRateDomain) getCny2USDRate() float64 {
	var cnyRateStr string
	var cnyRate = 7.3
	_ = d.cache.Get(keyUsdt2Cny, &cnyRateStr)
	if cnyRateStr != "" {
		cnyRate = tools.ToFloat64(cnyRateStr)
	}
	return cnyRate
}
