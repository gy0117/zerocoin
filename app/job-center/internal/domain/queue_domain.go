package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"job-center/internal/model"
	"zero-common/kafka"
)

const KLINE_1M = "kline_1m"
const topicBtcTransaction = "btc_transaction"

type QueueDomain struct {
	kafkaCli *kafka.KafkaClient
}

func NewQueueDomain(client *kafka.KafkaClient) *QueueDomain {
	return &QueueDomain{
		kafkaCli: client,
	}
}

func (d *QueueDomain) Send1mKline(ctx context.Context, data []string, symbol string) {
	kline := model.NewKline(data, "1m")
	bytes, err := json.Marshal(kline)
	if err != nil {
		logx.Error(err)
		return
	}

	item := kafka.KafkaData{
		Topic: KLINE_1M,
		Data:  bytes,
		Key:   []byte(symbol),
	}
	d.kafkaCli.Send(item)
	fmt.Println("job-center | Send1mKline success...")
}

func (d *QueueDomain) SendRecharge(value float64, time int64, addr string) {
	m := make(map[string]any)
	m["value"] = value
	m["address"] = addr
	m["time"] = time
	m["type"] = model.RECHARGE
	m["symbol"] = "BTC"
	marshal, err := json.Marshal(m)
	if err != nil {
		logx.Error(err)
		return
	}

	msg := kafka.KafkaData{
		Topic: topicBtcTransaction,
		Key:   []byte(addr),
		Data:  marshal,
	}
	d.kafkaCli.Send(msg)
	fmt.Println("job-center | SendRecharge success, and data: ", string(marshal))
}
