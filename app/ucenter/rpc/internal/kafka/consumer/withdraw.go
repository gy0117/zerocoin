package consumer

import (
	"common/kafka"
	"common/zerodb"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/model"
)

func ConsumeBtcWithdraw(kCli *kafka.KafkaClient, db *zerodb.ZeroDB) {
	withdrawDomain := domain.NewWithdrawDomain(db)
	for {
		kafkaData := kCli.Read()
		var record model.WithdrawRecord
		if err := json.Unmarshal(kafkaData.Data, &record); err != nil {
			logx.Error(err)
			continue
		}

		// 调用btc rpc进行提现
		err := withdrawDomain.Withdraw(record)
		if err != nil {
			logx.Error(err)
		} else {
			logx.Info("withdraw successfully")
		}
	}
}
