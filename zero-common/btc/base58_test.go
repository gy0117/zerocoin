package btc

import (
	"fmt"
	"testing"
)

func TestBaseEncode(t *testing.T) {
	encode := Base58Encode([]byte("zxcvbnmasdfghjkl"))
	fmt.Println("encode: ", string(encode))

	decode := Base58Decode(encode)
	fmt.Println("decode: ", string(decode))
}
