package btc

import (
	"encoding/json"
	"errors"
	"log"
	"zero-common/tools"
)

var apiUrl = "http://127.0.0.1:18332"
var auth = "Basic Yml0Y29pbjoxMjM0NTY="
var btcId = "zerocoin"

type ListUnspentInfoResp struct {
	Id     string              `json:"id"`
	Error  string              `json:"error"`
	Result []ListUnspentResult `json:"result"`
}

type ListUnspentResult struct {
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Confirmations int     `json:"confirmations"`
}

type Input struct {
	Txid string `json:"txid"`
	Vout int    `json:"vout"`
}
type CreateRawTransactionInfo struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

type SignRawTransactionWithWalletInfoResp struct {
	Id     string                             `json:"id"`
	Error  string                             `json:"error"`
	Result SignRawTransactionWithWalletResult `json:"result"`
}
type SignRawTransactionWithWalletResult struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

type SendRawTransactionInfoResp struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

// ListUnspent 查询账户地址的UTXO，拿到txid以及相关的信息
// https://developer.bitcoin.org/reference/rpc/listunspent.html
func ListUnspent(min, max int, addresses []string) ([]ListUnspentResult, error) {
	params := make(map[string]any)
	params["id"] = btcId
	params["method"] = "listunspent"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{min, max, addresses}

	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var resp ListUnspentInfoResp
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error == "" {
		return resp.Result, nil
	}
	return nil, errors.New(resp.Error)
}

// CreateRawTransaction 创建交易，返回hex string of the transaction
// https://developer.bitcoin.org/reference/rpc/createrawtransaction.html
func CreateRawTransaction(inputs []Input, values []map[string]any) (string, error) {
	params := make(map[string]any)
	params["id"] = btcId
	params["method"] = "createrawtransaction"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{inputs, values}
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var result CreateRawTransactionInfo
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return "", err
	}
	if result.Error == "" {
		return result.Result, nil
	}
	return "", errors.New(result.Error)
}

// SignRawTransactionWithWallet 将txid进行签名
// https://developer.bitcoin.org/reference/rpc/signrawtransactionwithwallet.html
func SignRawTransactionWithWallet(hexTxid string) (*SignRawTransactionWithWalletResult, error) {
	params := make(map[string]any)
	params["id"] = btcId
	params["method"] = "signrawtransactionwithwallet"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{hexTxid}
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var resp SignRawTransactionWithWalletInfoResp
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error == "" {
		return &resp.Result, nil
	}
	return nil, errors.New(resp.Error)
}

// SendRawTransaction 发送到btc网络，返回txid
// https://developer.bitcoin.org/reference/rpc/sendrawtransaction.html
func SendRawTransaction(signHex string) (string, error) {
	params := make(map[string]any)
	params["id"] = btcId
	params["method"] = "sendrawtransaction"
	params["jsonrpc"] = "1.0"
	params["params"] = []any{signHex, 0} //0代表任意手续费 前面创建交易的时候 一定要算好手续费
	headers := make(map[string]string)
	headers["Authorization"] = auth
	bytes, err := tools.PostWithHeader(apiUrl, params, headers, "")
	if err != nil {
		log.Println(err)
		return "", err
	}
	var resp SendRawTransactionInfoResp
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return "", err
	}

	if resp.Error == "" {
		return resp.Result, nil
	}
	return "", errors.New(resp.Error)
}
