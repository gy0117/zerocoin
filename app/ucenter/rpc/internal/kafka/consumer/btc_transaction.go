package consumer

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"ucenter-rpc/internal/domain"
	"zero-common/kafka"
	"zero-common/zerodb"
)

type BtcTransactionResult struct {
	Value   float64 `json:"value"`
	Time    int64   `json:"time"`
	Address string  `json:"address"`
	Type    string  `json:"type"`
	Symbol  string  `json:"symbol"`
}

// ConsumeBtcTransaction 从kafka中读取数据，存入数据库
func ConsumeBtcTransaction(kCli *kafka.KafkaClient, db *zerodb.ZeroDB) {
	transactionDomain := domain.NewTransactionDomain(db)

	for {
		kafkaData := kCli.Read()
		var bt BtcTransactionResult
		err := json.Unmarshal(kafkaData.Data, &bt)
		if err != nil {
			logx.Error(err)
			continue
		}
		err = transactionDomain.SaveRecharge(bt.Address, bt.Value, bt.Time, bt.Type, bt.Symbol)
		if err != nil {
			logx.Error(err)
			time.Sleep(time.Millisecond * 500)
			kCli.RepeatPut(kafkaData)
		}
	}
}
