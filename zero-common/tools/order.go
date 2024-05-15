package tools

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderId TODO-gy 可以使用雪花算法
func GenerateOrderId(prefix string) string {
	milli := time.Now().UnixMilli()
	val := rand.Intn(99999999)
	return fmt.Sprintf("%s%d%d", prefix, milli, val)
}
