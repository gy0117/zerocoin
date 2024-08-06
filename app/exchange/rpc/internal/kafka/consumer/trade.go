package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/trade"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
	"zero-common/kafka"
	"zero-common/zerodb"
)

// 消费订单的消息，拿到新创建的订单
// 1. 先实现买卖盘的逻辑，买入 卖出 一旦匹配完成，就成交了，成交的价格和数量，就会成为别人的参考，买卖盘实时的

const (
	topicExchangeOrderTrading           = "exchange_order_trading" // 将订单发送到撮合交易中
	topicExchangeOrderComplete          = "exchange_order_complete"
	topicUpdateWalletAfterOrderComplete = "update_wallet_after_order_complete"
)

type KafkaConsumer struct {
	cli      *kafka.KafkaClient
	factory  *trade.CoinTradeEngineFactory
	db       *zerodb.ZeroDB
	redisCli *redis.Redis
}

func NewKafkaConsumer(cli *kafka.KafkaClient, factory *trade.CoinTradeEngineFactory, db *zerodb.ZeroDB, redisCli *redis.Redis) *KafkaConsumer {
	return &KafkaConsumer{
		cli:      cli,
		factory:  factory,
		db:       db,
		redisCli: redisCli,
	}
}

func (kc *KafkaConsumer) Run() {
	kc.consumeTradingOrder()

	orderDomain := domain.NewOrderDomain(kc.db)
	kc.consumeCompleteOrder(orderDomain)
}

func (kc *KafkaConsumer) consumeCompleteOrder(orderDomain *domain.OrderDomain) {
	client := kc.cli.StartRead(topicExchangeOrderComplete)
	go kc.readCompleteOrder(client, orderDomain)
}

func (kc *KafkaConsumer) consumeTradingOrder() {
	client := kc.cli.StartRead(topicExchangeOrderTrading)
	go kc.readTradingOrder(client)
}

func (kc *KafkaConsumer) readTradingOrder(client *kafka.KafkaClient) {
	for {
		kafkaData := client.Read()
		logx.Info("[exchange-rpc] kafkaConsumer | readTradingOrder | topic: ", kafkaData.Topic, ", orderId: ", string(kafkaData.Key), ", kafkaData: ", string(kafkaData.Data))

		var exchangeOrder *model.ExchangeOrder
		err := json.Unmarshal(kafkaData.Data, &exchangeOrder)
		if err != nil {
			logx.Error("KafkaConsumer | Unmarshal err: ", err)
			continue
		}

		// kafka去重
		key := fmt.Sprintf("%s-%d-%s", kafkaData.Topic, exchangeOrder.UserId, exchangeOrder.OrderId)
		setnx, _ := kc.redisCli.Setnx(key, "1")
		if !setnx {
			logx.Error("KafkaConsumer | readTradingOrder | 消息重复，key：", key)
			continue
		}

		// 读取到订单数据，然后交给撮合交易引擎，由引擎处理
		coinTrade, ok := kc.factory.GetCoinTrade(exchangeOrder.Symbol)
		if !ok {
			logx.Error(errors.New("there is no corresponding matching trading engine available"))
			return
		}
		coinTrade.Trade(exchangeOrder)
	}
}

// 订单交易成功后，需要做两个处理
// 1. 更新数据库订单信息
// 2. 更新钱包
func (kc *KafkaConsumer) readCompleteOrder(client *kafka.KafkaClient, orderDomain *domain.OrderDomain) {
	for {
		kafkaData := client.Read()
		logx.Info("[exchange-rpc] | KafkaConsumer | readCompleteOrder | topic: ", kafkaData.Topic, ", orderId: ", string(kafkaData.Key), ", kafkaData: ", string(kafkaData.Data))

		var exchangeOrder *model.ExchangeOrder
		err := json.Unmarshal(kafkaData.Data, &exchangeOrder)
		if err != nil {
			logx.Error("[exchange-rpc] | KafkaConsumer | Unmarshal err: ", err)
			continue
		}

		// 1. 更新订单信息
		err = orderDomain.UpdateOrderComplete(context.Background(), exchangeOrder)
		if err != nil {
			logx.Error(err)
			client.RepeatPut(kafkaData)
			time.Sleep(time.Millisecond * 250)
			continue
		}

		// 2. 通知钱包更新（kafka）
		for {
			kafkaData.Topic = topicUpdateWalletAfterOrderComplete
			err2 := client.SendSync(kafkaData)
			if err2 != nil {
				logx.Error(fmt.Sprintf("通知钱包更新失败，err: %s\n", err2.Error()))
				time.Sleep(time.Millisecond * 250)
				continue
			}
			logx.Info("[exchange-rpc] | KafkaConsumer ｜ 发送 update_wallet_after_order_complete 消息, orderId: " + exchangeOrder.OrderId)
			break
		}
	}
}
