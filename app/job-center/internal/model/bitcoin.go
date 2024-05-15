package model

type MiningInfoResp struct {
	Id     string         `json:"id"`
	Error  string         `json:"error"`
	Result MiningInfoData `json:"result"`
}

type BlockHashResp struct {
	Id     string `json:"id"`
	Error  string `json:"error"`
	Result string `json:"result"`
}

type BlockResp struct {
	Id     string    `json:"id"`
	Error  string    `json:"error"`
	Result BlockData `json:"result"`
}

type RawTransactionResp struct {
	Id     string             `json:"id"`
	Error  string             `json:"error"`
	Result RawTransactionData `json:"result"`
}

type MiningInfoData struct {
	Blocks        int     `json:"blocks"`        // 区块数量
	Difficulty    float64 `json:"difficulty"`    // 出块难度
	Networkhashps float64 `json:"networkhashps"` // 全网hash生成速率
	Pooledtx      int     `json:"pooledtx"`      // 内存交易池中的交易数量
	Chain         string  `json:"chain"`         // 当前所连接网络
	Warnings      string  `json:"warnings"`
}

type BlockData struct {
	Hash string `json:"hash"` // 区块hash
	//Confirmations     int      `json:"confirmations"` // 确认书
	//Height            int64    `json:"height"`        // 区块高度
	//Version           int64    `json:"version"`       // 版本
	//VersionHex        string   `json:"versionHex"`    // 16进制表示的版本
	//MerkleRoot        string   `json:"merkleroot"`    // 区块的默克尔树根
	Time int64 `json:"time"` // 区块创建时间戳
	//MedianTime        int64    `json:"mediantime"`    // 区块中值时间戳
	//Nonce             int64    `json:"nonce"`
	//Bits              string   `json:"bits"`
	//Difficulty        float64  `json:"difficulty"` // 难度
	//ChainWork         string   `json:"chainwork"`
	//NTx               int      `json:"nTx"`
	//PreviousBlockHash string   `json:"previousBblockhash"` // 前一区块的哈希
	//NextBlockHash     string   `json:"nextblockhash"`      // 下一区块的哈希
	//Strippedsize      int64    `json:"strippedsize"`
	//Size              int64    `json:"size"`
	//Weight            int64    `json:"weight"`
	Tx []string `json:"tx"` // 交易id数组
}

type RawTransactionData struct {
	TxId      string `json:"txid"`
	Hash      string `json:"hash"`
	Version   int    `json:"version"`
	Size      int    `json:"size"`
	Vsize     int    `json:"vsize"`
	Weight    int    `json:"weight"`
	LockTime  int64  `json:"locktime"`
	Vin       []Vin  `json:"vin"`
	Vout      []Vout `json:"vout"`
	Hex       string `json:"hex"`
	BlockHash string `json:"blockhash"`
	//Confirmations int64  `json:"confirmations"`
	Time      int64 `json:"time"`
	BlockTime int64 `json:"blocktime"`
}

type Vin struct {
	TxId        string            `json:"txid"`
	Vout        int               `json:"vout"`
	ScriptSig   map[string]string `json:"scriptSig"`
	Txinwitness []string          `json:"txinwitness"`
	Sequence    int64             `json:"sequence"`
}

//type ScriptSig struct {
//	Asm string `json:"asm"`
//	Hex string `json:"hex"`
//}

type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
