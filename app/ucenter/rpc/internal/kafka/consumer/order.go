package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"grpc-common/exchange/eclient"
	"grpc-common/exchange/types/order"
	"time"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/model"
	"zero-common/kafka"
	"zero-common/operate"
	"zero-common/zerodb"
)

const (
	DirectionBuy  = "BUY"
	DirectionSell = "SELL"

	DirectionBuyInt  = 0
	DirectionSellInt = 1
)

// 钱包的钱冻结成功后，需要修改钱包状态
const topicExchangeOrderInitComplete = "exchange_order_init_complete"

type OrderAdd struct {
	UserId     int64   `json:"userId"`
	OrderId    string  `json:"orderId"`
	Money      float64 `json:"money"`
	Symbol     string  `json:"symbol"`
	Direction  int     `json:"direction"`
	BaseSymbol string  `json:"baseSymbol"`
	CoinSymbol string  `json:"coinSymbol"`
}

func UpdateWalletAfterOrderComplete(kCli *kafka.KafkaClient, redisCli *redis.Redis, db *zerodb.ZeroDB) {
	walletDomain := domain.NewWalletDomain(db)

	for {
		kafkaData := kCli.Read()
		var exchangeOrder *model.ExchangeOrder
		err := json.Unmarshal(kafkaData.Data, &exchangeOrder)
		if err != nil {
			logx.Error(fmt.Sprintf("[ucenter-rpc] ｜ UpdateWalletAfterOrderComplete ｜ Unmarshal err: %s\n", err.Error()))
			time.Sleep(time.Millisecond * 250)
			continue
		}

		// check
		if exchangeOrder == nil || exchangeOrder.Status != model.StatusCompleted {
			continue
		}
		logx.Info("[ucenter-rpc] | 收到 update_wallet_after_order_complete 消息 | orderId: " + exchangeOrder.OrderId)

		// @gy 重点
		// user_id的row 只能同时由一个人操作，使用gozero自带的锁
		redisLock := redis.NewRedisLock(redisCli, fmt.Sprintf("update_wallet_after_order_complete::%d", exchangeOrder.UserId))
		acquire, err := redisLock.Acquire()
		if err != nil {
			logx.Error("[ucenter-rpc] | 有人已经拿到锁，正在处理中, err: ", err)
			continue
		}
		if acquire {
			err2 := updateWalletAfterOrderCompleteInner(exchangeOrder, walletDomain)
			if err2 != nil {
				kCli.RepeatPut(kafkaData)
				time.Sleep(time.Millisecond * 250)
				_, _ = redisLock.Release()
				continue
			}
			_, _ = redisLock.Release()
		}
	}
}

func updateWalletAfterOrderCompleteInner(exchangeOrder *model.ExchangeOrder, walletDomain *domain.WalletDomain) error {
	// 1. 买单
	// 1.1 市价买：冻结的钱是amount USDT；扣的钱是order.turnover；还回去的钱是amount - order.turnover
	// 1.2 限价买：冻结的钱是order.price * amount；成交的钱是turnover；还回去的钱是order.price * amount - turnover

	// 2. 卖单
	// 2.1 不管限价还是市价，卖的都是BTC，解冻的钱是amount，得到的钱是order.turnover

	ctx := context.Background()
	// 先查钱包
	coinWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, exchangeOrder.UserId, exchangeOrder.CoinSymbol)
	if err != nil {
		logx.Error("[ucenter-rpc] | updateMemberWalletAfterOrderComplete | buy | FindWalletByMemIdAndCoinName coinWallet, err: ", err)
		return err
	}
	baseWallet, err := walletDomain.FindWalletByMemIdAndCoinName(ctx, exchangeOrder.UserId, exchangeOrder.BaseSymbol)
	if err != nil {
		logx.Error("[ucenter-rpc] | updateMemberWalletAfterOrderComplete | buy | FindWalletByMemIdAndCoinName baseWallet, err: ", err)
		return err
	}

	if exchangeOrder.Direction == model.DirectionBuyInt {
		if exchangeOrder.Type == model.MarketPriceInt {
			// 以市价买，买入量就是花的钱。花usdt，买btc；例如：花1000usdt，买btc
			baseWallet.FrozenBalance = operate.SubFloor(baseWallet.FrozenBalance, exchangeOrder.Amount, 8)
			baseWallet.Balance = operate.AddFloor(baseWallet.Balance, operate.SubFloor(exchangeOrder.Amount, exchangeOrder.Turnover, 8), 8)
			coinWallet.Balance = operate.AddFloor(coinWallet.Balance, exchangeOrder.TradedAmount, 8)
		} else {
			// 以限价买
			x := operate.MulFloor(exchangeOrder.Price, exchangeOrder.Amount, 8)
			baseWallet.FrozenBalance = operate.SubFloor(baseWallet.FrozenBalance, x, 8)
			baseWallet.Balance = operate.AddFloor(baseWallet.Balance, operate.SubFloor(x, exchangeOrder.Turnover, 8), 8)
			coinWallet.Balance = operate.AddFloor(coinWallet.Balance, exchangeOrder.TradedAmount, 8)
		}
	} else {
		//卖 不管是市价还是限价 都是卖的 BTC  解冻amount 得到的钱是 order.turnover
		// 打算卖10个btc，这一单成交2个btc，此时就还冻结8个btc；usdt就增加了对应的成交额
		coinWallet.FrozenBalance = operate.SubFloor(coinWallet.FrozenBalance, exchangeOrder.Amount, 8)
		baseWallet.Balance = operate.AddFloor(baseWallet.Balance, exchangeOrder.Turnover, 8)
	}

	// 更新钱包信息
	err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
	if err != nil {
		logx.Error("[ucenter-rpc] | updateMemberWalletAfterOrderComplete | buy | UpdateWalletCoinAndBase, err: ", err)
		return err
	}
	logx.Info("[ucenter-rpc] | 更新钱包成功 | orderId: " + exchangeOrder.OrderId)
	return nil
}

// ExchangeOrderAdd 订单创建成功后，需要冻结钱包的钱
func ExchangeOrderAdd(kCli *kafka.KafkaClient, orderRpc eclient.Order, db *zerodb.ZeroDB, redisCli *redis.Redis) {
	for {
		// 1. kafka 消费数据
		data := kCli.Read()
		orderId := string(data.Key)
		var orderAdd OrderAdd
		_ = json.Unmarshal(data.Data, &orderAdd)
		logx.Info("【ucenter-rpc】received message for order creation，orderId: ", orderId)

		if orderId != orderAdd.OrderId {
			logx.Error("order inconsistency")
			continue
		}

		// 2. 根据orderId，查询订单状态，如果订单交易中，说明已经被处理了
		ctx := context.Background()
		exchangeOrder, err := orderRpc.FindByOrderId(ctx, &order.OrderReq{
			OrderId: orderId,
		})
		if err != nil {
			logx.Error(err)
			cancelOrder(ctx, kCli, orderRpc, orderId, data)
			continue
		}
		if exchangeOrder == nil {
			logx.Errorf("orderId %s order does not exist", orderId)
			continue
		}
		if exchangeOrder.Status != model.StatusInit {
			logx.Errorf("orderId %s has already been processed", orderId)
			continue
		}

		lock := redis.NewRedisLock(redisCli, fmt.Sprintf("add_order_freeze_wallet::%d", exchangeOrder.UserId))
		acquire, err := lock.Acquire()
		if err != nil {
			logx.Error("a process has already acquired the lock, err: ", err)
			continue
		}
		if acquire {
			walletDomain := domain.NewWalletDomain(db)
			// 3. 冻结钱包
			if orderAdd.Direction == DirectionBuyInt {
				// 买入，使用USDT
				err = walletDomain.Freeze(ctx, orderAdd.UserId, orderAdd.Money, orderAdd.BaseSymbol)
			} else {
				// 卖出，卖出btc等
				err = walletDomain.Freeze(ctx, orderAdd.UserId, orderAdd.Money, orderAdd.CoinSymbol)
			}
			if err != nil {
				logx.Errorf("failed to freeze wallet, order direction: %v, err: %v", orderAdd.Direction, err)
				cancelOrder(ctx, kCli, orderRpc, orderId, data)
				continue
			}

			// 4. 冻结成功后，需要修改订单状态，改成trading
			for {
				m := make(map[string]any)
				m["userId"] = orderAdd.UserId
				m["orderId"] = orderId
				marshal, _ := json.Marshal(m)
				kData := kafka.KafkaData{
					Topic: topicExchangeOrderInitComplete,
					Key:   []byte(orderId),
					Data:  marshal,
				}

				err := kCli.SendSync(kData)
				if err != nil {
					logx.Error(err)
					time.Sleep(time.Millisecond * 500)
					continue
				}
				logx.Info("【ucenter-rpc】冻结成功后，发送修改订单状态消息成功，orderId: ", orderId)
				break
			}
			if _, err := lock.Release(); err != nil {
				logx.Error("lock release, err： ", err)
			}
		}
	}
}

// 重新消费
func cancelOrder(ctx context.Context, kCli *kafka.KafkaClient, orderRpc eclient.Order, orderId string, kData kafka.KafkaData) {
	// 1. 取消订单
	_, err := orderRpc.CancelOrder(ctx, &order.OrderReq{
		OrderId: orderId,
	})
	if err != nil {
		logx.Error("failed to cancel order, err: ", err)
		// 2. 订单消息重新添加到kafka中
		kCli.RepeatPut(kData)
	}
}
