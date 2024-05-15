package model

const (
	RECHARGE = "RECHARGE"
	WITHDRAW = "WITHDRAW"
)

const TableBtcTransaction = "btc_transaction"

type BtcTransaction struct {
	TxId      string  `bson:"txId"`
	Time      int64   `bson:"time"`
	Value     float64 `bson:"value"`
	BlockHash string  `bson:"blockhash"`
	Address   string  `bson:"address"`
	Type      string  `bson:"type"` // RECHARGE 充值 WITHDRAW 提现
}
