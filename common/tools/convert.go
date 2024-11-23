package tools

import (
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func ToInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logx.Error(err)
		return 0
	}
	return v
}

func ToFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logx.Error(err)
		return 0
	}
	return v
}
