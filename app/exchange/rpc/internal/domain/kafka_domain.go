package domain

import (
	"context"
	"encoding/json"
	"errors"
	"exchange-rpc/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"time"
	"zero-common/kafka"
)

const (
	topicExchangeOrderInitComplete = "exchange_order_init_complete" // 钱包的钱冻结成功后，需要修改钱包状态
	topicExchangeOrderTrading      = "exchange_order_trading"       // 将订单发送到撮合交易中
)

type OrderResult struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

type KafkaDomain struct {
	cli         *kafka.KafkaClient
	orderDomain *OrderDomain
}

func NewKafkaDomain(cli *kafka.KafkaClient, order *OrderDomain) *KafkaDomain {
	k := &KafkaDomain{
		cli:         cli,
		orderDomain: order,
	}
	go k.WaitAddOrderResult()
	return k
}

func (d *KafkaDomain) SendOrderAdd(topic string, userId int64, orderId string, money float64, symbol string, direction int, baseSymbol string, coinSymbol string) error {
	m := make(map[string]any)
	m["topic"] = topic
	m["userId"] = userId
	m["orderId"] = orderId
	m["money"] = money
	m["symbol"] = symbol
	m["direction"] = direction
	m["baseSymbol"] = baseSymbol
	m["coinSymbol"] = coinSymbol
	bytes, _ := json.Marshal(m)

	kData := kafka.KafkaData{
		Topic: topic,
		Key:   []byte(orderId),
		Data:  bytes,
	}

	if err := d.cli.SendSync(kData); err != nil {
		logx.Error(err)
		return errors.New("the order message failed to send to Kafka")
	}
	logx.Info("【exchange-rpc】order created, message sent successfully，orderId: ", orderId)
	return nil
}

// WaitAddOrderResult 接收【冻结之后，修改订单状态的消息】
// 必须确保执行成功
func (d *KafkaDomain) WaitAddOrderResult() {
	logx.Info("【exchange-rpc】waitAddOrderResult | start modify order status")
	client := d.cli.StartRead(topicExchangeOrderInitComplete)
	for {
		kData := client.Read()
		//logx.Info("【exchange-rpc】读取exchange_order_init_complete 消息成功: " + string(kData.Key))
		var orderResult OrderResult
		if err := json.Unmarshal(kData.Data, &orderResult); err != nil {
			logx.Error("开始修改订单状态 | unmarshal err：", err)
			continue
		}
		logx.Infof("【exchange-rpc】冻结成功后，开始修改订单状态 | orderResult：%+v\n", orderResult)
		ctx := context.Background()
		// 查询订单
		exchangeOrder, err := d.orderDomain.FindByOrderId(ctx, &order.OrderReq{
			OrderId: orderResult.OrderId,
		})
		if err != nil {
			logx.Error("【exchange-rpc】冻结成功后，开始修改订单状态 | FindByOrderId err：", err)
			//client.RepeatPut(kData)
			continue
		}
		if exchangeOrder == nil {
			logx.Errorf("开始修改订单状态 | %s 订单不存在", orderResult.OrderId)
			continue
		}
		if exchangeOrder.Status != model.OrderStatus_StatusInit {
			logx.Errorf("开始修改订单状态 | %s 订单已经被处理过了", orderResult.OrderId)
			continue
		}
		logx.Infof("【exchange-rpc】冻结成功后，开始修改订单状态 | FindByOrderId exchangeOrder：%+v\n", exchangeOrder)

		// 将订单状态改成"交易中"
		err = d.orderDomain.UpdateOrderStatus(ctx, orderResult.OrderId, model.OrderStatus_Trading)
		if err != nil {
			logx.Error("开始修改订单状态 | UpdateOrderStatus err：", err)
			client.RepeatPut(kData)
			continue
		}

		// 需要发送消息到kafka，订单需要加入到撮合交易当中
		// 如果没有撮合交易成功，加入到撮合交易等待队列，继续等待完成撮合
		exchangeOrder.Status = model.OrderStatus_Trading
		for {
			marshal, _ := json.Marshal(exchangeOrder)
			logx.Info("【exchange-rpc】订单状态修改成功 | 发送到kafka：", string(marshal))

			orderData := kafka.KafkaData{
				Topic: topicExchangeOrderTrading,
				Key:   []byte(exchangeOrder.OrderId),
				Data:  marshal,
			}
			err = client.SendSync(orderData)
			if err != nil {
				logx.Error(err)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			logx.Info("【exchange-rpc】订单状态修改成功 | 发送订单数据到撮合交易，orderId: ", string(orderData.Key))
			break
		}
	}
}
