package domain

import (
	"strings"
	"zero-common/zerodb"
)

type ExchangeRateDomain struct {
}

func NewExchangeRateDomain(db *zerodb.ZeroDB) *ExchangeRateDomain {
	return &ExchangeRateDomain{}
}

// UsdRate TODO 从redis中查询，在job-center中写一个定时任务查询实时汇率，存入redis
func (d *ExchangeRateDomain) UsdRate(unit string) float64 {
	unit = strings.ToUpper(unit)
	if "CNY" == unit {
		return 7.3
	}
	if "JPY" == unit {
		return 137
	}
	return 0
}
