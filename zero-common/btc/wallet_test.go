package btc

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		panic(err)
	}
	address := wallet.GetAddress()
	fmt.Println(string(address))

	// https://www.blockchain.com/explorer/addresses/btc/ + string(address)，可查询钱包
}
