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
	//address := wallet.GenerateBitcoinAddress()
	//fmt.Println(string(address))

	// mm4W2wPtVCoWCUvWVwk4v7rrGy7UqPcSsX （已经申请）
	// mpqS3daJAHZCLfW6owByDu17qvn62u41Me
	// mqqQhmE2ATmGGzdF1KaeFYfYiGYGSgwNNA
	// mwJ4rU484frfaC1mXfqdJFMypUfikcywkK
	// n16y9qGbP4hc9hsV7pf7H3u1q1M7EDRaSH
	// n16y9qGbP4hc9hsV7pf7H3u1q1M7EDRaSH
	// n16y9qGbP4hc9hsV7pf7H3u1q1M7EDRaSH
	testAddress := wallet.GenerateBitcoinTestAddress()
	fmt.Println(string(testAddress))

	privateKey := wallet.GenerateBitcoinPrivateKey()
	fmt.Println(privateKey)

	// https://www.blockchain.com/explorer/addresses/btc/ + string(address)，可查询钱包
}

// 查看测试网络bitcoin
// https://blockstream.info/testnet/address/mm4W2wPtVCoWCUvWVwk4v7rrGy7UqPcSsX
// https://live.blockcypher.com/btc-testnet/address/mm4W2wPtVCoWCUvWVwk4v7rrGy7UqPcSsX/
