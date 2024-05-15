package tools

import (
	"fmt"
	"math/rand"
)

func GenerateVerifyCode() string {
	num := rand.Intn(9999)
	if num < 1000 {
		num += 1000
	}
	return fmt.Sprintf("%d", num)
}
