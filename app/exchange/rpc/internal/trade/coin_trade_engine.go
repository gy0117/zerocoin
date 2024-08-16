package trade

import (
	"context"
	"encoding/json"
	"errors"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/trade/plate"
	"exchange-rpc/internal/trade/queue"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"zero-common/kafka"
	"zero-common/operate"
	"zero-common/zerodb"
)

const (
	topicExchangeOrderTradePlate = "exchange_order_trade_plate" // 发送盘口信息到kafka，在market服务中读取kafka消息，然后推送到前端
	topicExchangeOrderComplete   = "exchange_order_complete"
)

type CoinTradeEngine struct {
	symbol string
	kCli   *kafka.KafkaClient

	marketTradeOperate *MarketTradeOperate // 市价操作对象
	limitTradeOperate  *LimitTradeOperate  // 限价操作对象

	buyTradePlate  *plate.TradePlate // 买盘
	sellTradePlate *plate.TradePlate // 卖盘
}

// Trade 撮合交易核心代码
// 1. 当订单进来之后，我们判断 buy 还是 sell
// 2. 确定 市价 还是 限价
// 3. buy：和sell队列进行匹配；sell：和buy队列进行匹配
// 4. 里面会有成交，还未交易的，放入买卖盘
// 5. 订单就会更新，订单的状态要变，冻结的金额扣除等等
func (engine *CoinTradeEngine) Trade(order *model.ExchangeOrder) {
	logx.Info("[exchange-rpc] Trade | start buy&sell trade plate")

	var limitPriceList *queue.LimitPriceQueue
	var marketPriceList queue.TradeTimeQueue

	if order.Direction == model.DirectionBuyInt {
		// 如果是买订单，则需要跟卖出队列进行匹配
		limitPriceList = engine.limitTradeOperate.sellLimitQueue
		marketPriceList = engine.marketTradeOperate.sellMarketQueue
	} else if order.Direction == model.DirectionSellInt {
		// 如果是卖订单，则需要跟买入队列进行匹配
		limitPriceList = engine.limitTradeOperate.buyLimitQueue
		marketPriceList = engine.marketTradeOperate.buyMarketQueue
	}

	if order.Type == model.MarketPriceInt {
		// 市价单，市价订单和限价订单进行匹配
		engine.matchMarketPriceWithLimitPrice(limitPriceList, order)
	}
	if order.Type == model.LimitPriceInt {
		// 限价单，先与限价单进行成交；如果未成交，继续与市价单进行成交
		engine.matchLimitPriceWithLimitPrice(limitPriceList, order)

		if order.Status == model.OrderStatus_Trading {
			engine.matchLimitPriceWithMarketPrice(marketPriceList, order)
		}

		// 如果还没成交，则加入到买卖盘
		if order.Status == model.OrderStatus_Trading {
			engine.addLimitPriceQueue(order)
			if order.Direction == model.DirectionBuyInt {
				engine.sendTradePlateData(engine.buyTradePlate)
			} else {
				engine.sendTradePlateData(engine.sellTradePlate)
			}
		}
	}
}

func (engine *CoinTradeEngine) init(db *zerodb.ZeroDB) {
	engine.marketTradeOperate = NewMarketTradeOperate()
	engine.limitTradeOperate = NewLimitTradeOperate()

	engine.buyTradePlate = plate.NewTradePlate(engine.symbol, plate.BUY)
	engine.sellTradePlate = plate.NewTradePlate(engine.symbol, plate.SELL)

	engine.initData(db)
}

// 将盘口信息发送到kafka，在market服务中，读取kafka消息，然后推到前端
func (engine *CoinTradeEngine) sendTradePlateData(plate *plate.TradePlate) {
	plateResult := plate.NewResult()
	marshal, _ := json.Marshal(plateResult)

	kData := kafka.KafkaData{
		Topic: topicExchangeOrderTradePlate,
		Key:   []byte(plate.Symbol),
		Data:  marshal,
	}
	err := engine.kCli.SendSync(kData)
	if err != nil {
		logx.Error(err)
		return
	}
}

// 买卖盘的数据展示
// 查exchange_order表
func (engine *CoinTradeEngine) initData(db *zerodb.ZeroDB) {
	orderDomain := domain.NewOrderDomain(db)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	exchangeOrders, err := orderDomain.FindTradingOrderBySymbol(ctx, engine.symbol)
	if err != nil {
		logx.Error(err)
		return
	}
	if exchangeOrders == nil {
		logx.Error(errors.New("no records found"))
		return
	}

	// 1. 区分市价or限价，市价不需要进入买卖盘，限价单才会进入买卖盘
	// 2. 区分买or卖
	for _, o := range exchangeOrders {
		if o.Type == model.MarketPriceInt {

			if o.Direction == model.DirectionBuyInt {
				// 市价 && 买入 订单
				engine.marketTradeOperate.JoinInBuyMarketQueueSafe(o)
			} else if o.Direction == model.DirectionSellInt {
				// 市价 && 卖出 订单
				engine.marketTradeOperate.JoinInSellMarketQueueSafe(o)
			}

		} else if o.Type == model.LimitPriceInt {

			if o.Direction == model.DirectionBuyInt {
				// 限价 && 买入 订单
				// 1. 加入到限价队列； 2. 加入到买盘
				engine.limitTradeOperate.buyLimitQueue.Lock()
				engine.limitTradeOperate.JoinInBuyLimitQueue(o)
				engine.buyTradePlate.Add(o)
				engine.limitTradeOperate.buyLimitQueue.Unlock()

			} else if o.Direction == model.DirectionSellInt {
				// 限价 && 卖出 订单
				// 1. 加入到限价队列； 2. 加入到卖盘
				engine.limitTradeOperate.sellLimitQueue.Lock()
				engine.limitTradeOperate.JoinInSellLimitQueue(o)
				engine.sellTradePlate.Add(o)
				engine.limitTradeOperate.sellLimitQueue.Unlock()
			}
		}
	}
	// 市价，按照订单时间排序
	engine.marketTradeOperate.Sort()

	// 限价，如果是买入订单队列，则按照价格从高到底排序（买入的价格越高，则越容易成交）
	// 如果是卖出订单队列，则按照价格从低到高排序（卖出的价格越低，则越容易成交）
	engine.limitTradeOperate.Sort()

	// 发送盘口信息到kafka，然后在market服务中，读取kafka消息，展示到前端
	if len(exchangeOrders) > 0 {
		engine.sendTradePlateData(engine.buyTradePlate)
		engine.sendTradePlateData(engine.sellTradePlate)
	}
}

// 使用市价队列，匹配限价单
func (engine *CoinTradeEngine) matchLimitPriceWithMarketPrice(list queue.TradeTimeQueue, curOrder *model.ExchangeOrder) {
	delOrders := make([]string, 0)

	for _, matchOrder := range list {
		if matchOrder.UserId == curOrder.UserId {
			continue
		}

		matchAmount := operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
		if matchAmount <= 0 {
			continue
		}
		unTradedAmount := operate.SubFloor(curOrder.Amount, curOrder.TradedAmount, 8)

		// 跟需求有关，我这里以限价为标准，符合条件的都以限价成交
		price := curOrder.Price // 限价价格

		if matchAmount >= unTradedAmount {
			// 能满足订单，直接交易
			matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, unTradedAmount, 8)

			curTurnover := operate.MulFloor(price, unTradedAmount, 8)
			matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)

			// 如果matchOrder用完
			if operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
				matchOrder.Status = model.OrderStatus_Completed
				// 删除 matchOrder
				delOrders = append(delOrders, matchOrder.OrderId)
			}

			curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, unTradedAmount, 8)
			curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)
			curOrder.Status = model.OrderStatus_Completed

			// 市价没有买卖盘
			break
		} else {
			// matchOrder的amount不够，先有多少用多少
			matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, matchAmount, 8)

			curTurnover := operate.MulFloor(price, matchAmount, 8)
			matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)
			matchOrder.Status = model.OrderStatus_Completed
			delOrders = append(delOrders, matchOrder.OrderId)

			curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, matchAmount, 8)
			curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)

			// 市价没有买卖盘
			// 继续下一轮匹配
			continue
		}
	}
	for _, orderId := range delOrders {
		for i, matchOrder := range list {
			if matchOrder.OrderId == orderId {
				list = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
}

// 使用限价队列，匹配限价单
func (engine *CoinTradeEngine) matchLimitPriceWithLimitPrice(list *queue.LimitPriceQueue, curOrder *model.ExchangeOrder) {
	list.Lock()
	defer list.Unlock()

	delOrders := make([]string, 0)
	buyNotify := false
	sellNotify := false
	var completeOrders []*model.ExchangeOrder

	for _, v := range list.List {
		for _, matchOrder := range v.List {
			if matchOrder.UserId == curOrder.UserId {
				continue
			}

			// 买入订单，如果买的价格比卖的价格还低，则无法成交
			// sellLimitList是根据价格从低到高排序的
			// 如果买入订单的价格是100，卖队列的订单价格为110，120，130，...
			if curOrder.Direction == model.DirectionBuyInt && curOrder.Price < matchOrder.Price {
				break
			}

			// 卖出订单，如果卖的价格比买的价格还高，则无法成交
			// buyLimitList是根据价格从高到低排序的
			// 如果卖出订单的价格是100，买队列的订单价格为90，80，70，...
			if curOrder.Direction == model.DirectionSellInt && curOrder.Price > matchOrder.Price {
				break
			}

			//price := matchOrder.Price
			matchAmount := operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}
			unTradedAmount := operate.SubFloor(curOrder.Amount, curOrder.TradedAmount, 8)

			if matchAmount >= unTradedAmount {
				// 能满足订单，直接交易
				matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, unTradedAmount, 8)

				curTurnover := operate.MulFloor(matchOrder.Price, unTradedAmount, 8)
				matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)

				// 如果matchOrder用完
				if operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.OrderStatus_Completed

					// 删除 matchOrder
					delOrders = append(delOrders, matchOrder.OrderId)
					completeOrders = append(completeOrders, matchOrder)
				}

				curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, unTradedAmount, 8)
				curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)
				curOrder.Status = model.OrderStatus_Completed
				completeOrders = append(completeOrders, curOrder)

				// 更新买卖盘
				if matchOrder.Direction == model.DirectionBuyInt {
					engine.buyTradePlate.Remove(matchOrder, unTradedAmount)
					buyNotify = true
				} else {
					engine.sellTradePlate.Remove(matchOrder, unTradedAmount)
					sellNotify = true
				}
				break
			} else {
				// matchOrder的amount不够，先有多少匹配多少
				matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, matchAmount, 8)

				curTurnover := operate.MulFloor(matchOrder.Price, matchAmount, 8)
				matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)
				matchOrder.Status = model.OrderStatus_Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				completeOrders = append(completeOrders, matchOrder)

				curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, matchAmount, 8)
				curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)

				// 更新买卖盘
				if matchOrder.Direction == model.DirectionBuyInt {
					engine.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					engine.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				// 继续下一轮匹配
				continue
			}
		}
	}
	// 删除order
	for _, orderId := range delOrders {
		for _, v := range list.List {
			for i, matchOrder := range v.List {
				if matchOrder.OrderId == orderId {
					v.List = append(v.List[:i], v.List[i+1:]...)
					break
				}
			}
		}
	}
	if buyNotify {
		engine.sendTradePlateData(engine.buyTradePlate)
	}
	if sellNotify {
		engine.sendTradePlateData(engine.sellTradePlate)
	}
	for _, v := range completeOrders {
		engine.sendCompleteOrder(v)
	}
}

// 使用限价队列，匹配市价单
func (engine *CoinTradeEngine) matchMarketPriceWithLimitPrice(list *queue.LimitPriceQueue, curOrder *model.ExchangeOrder) {
	list.Lock()
	defer list.Unlock()

	delOrders := make([]string, 0)
	buyNotify := false
	sellNotify := false

	for _, v := range list.List {
		for _, matchOrder := range v.List {
			// 自己的订单不处理
			if matchOrder.UserId == curOrder.UserId {
				continue
			}

			// 计算可交易的数量
			matchAmount := operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}
			//price := matchOrder.Price
			unTradedAmount := operate.SubFloor(curOrder.Amount, curOrder.TradedAmount, 8)

			// 如果是市价买，amount就是USDT的数量，要计算买多少BTC，要根据match的price进行计算
			// 举例子：卖出10btc，买入10000u的btc，这里就应该计算10000u可以买多少个btc
			if curOrder.Direction == model.DirectionBuyInt {
				// 总钱数 - 成交额，即还有这么多钱，要买入多少个btc呢？ 剩余钱数 / btc价格
				tmp := operate.SubFloor(curOrder.Amount, curOrder.Turnover, 8)
				unTradedAmount = operate.DivFloor(tmp, matchOrder.Price, 8)
			}
			if matchAmount >= unTradedAmount {
				// 能满足订单，直接交易
				matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, unTradedAmount, 8)

				curTurnover := operate.MulFloor(matchOrder.Price, unTradedAmount, 8)
				matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)

				// 如果matchOrder用完
				if operate.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.OrderStatus_Completed
					// 删除 matchOrder
					delOrders = append(delOrders, matchOrder.OrderId)
				}

				// 处理当前订单
				curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, unTradedAmount, 8)
				curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)
				curOrder.Status = model.OrderStatus_Completed

				// 更新买卖盘
				if matchOrder.Direction == model.DirectionBuyInt {
					engine.buyTradePlate.Remove(matchOrder, unTradedAmount)
					buyNotify = true
				} else {
					engine.sellTradePlate.Remove(matchOrder, unTradedAmount)
					sellNotify = true
				}
				break
			} else {
				// matchOrder的amount不够，那就先把可以交易的份额都成交了
				matchOrder.TradedAmount = operate.AddFloor(matchOrder.TradedAmount, matchAmount, 8)

				curTurnover := operate.MulFloor(matchOrder.Price, matchAmount, 8)
				matchOrder.Turnover = operate.AddFloor(matchOrder.Turnover, curTurnover, 8)
				matchOrder.Status = model.OrderStatus_Completed

				delOrders = append(delOrders, matchOrder.OrderId)

				// 更新当前订单的交易量和交易额
				curOrder.TradedAmount = operate.AddFloor(curOrder.TradedAmount, matchAmount, 8)
				curOrder.Turnover = operate.AddFloor(curOrder.Turnover, curTurnover, 8)

				// 更新买卖盘
				if matchOrder.Direction == model.DirectionBuyInt {
					engine.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					engine.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				// 当前订单还可以交易，继续下一轮匹配
				continue
			}
		}
	}

	// 删除order
	for _, orderId := range delOrders {
		for _, v := range list.List {
			for i, matchOrder := range v.List {
				if matchOrder.OrderId == orderId {
					v.List = append(v.List[:i], v.List[i+1:]...)
					break
				}
			}
		}
	}

	// 判断订单是否完成，如果没有完成，放入队列
	if curOrder.Status == model.OrderStatus_Trading {
		// 说明没有完成
		engine.addMarketPriceQueue(curOrder)
	}

	// 买卖盘更新
	if buyNotify {
		engine.sendTradePlateData(engine.buyTradePlate)
	}
	if sellNotify {
		engine.sendTradePlateData(engine.sellTradePlate)
	}
}

func (engine *CoinTradeEngine) addMarketPriceQueue(order *model.ExchangeOrder) {
	engine.marketTradeOperate.JoinInMarketPriceQueue(order)
}

func (engine *CoinTradeEngine) addLimitPriceQueue(order *model.ExchangeOrder) {
	if order.Type != model.LimitPriceInt {
		return
	}
	if order.Direction == model.DirectionBuyInt {
		engine.limitTradeOperate.buyLimitQueue.Lock()
		engine.limitTradeOperate.JoinInBuyLimitQueue(order)
		// 加入到买盘
		engine.buyTradePlate.Add(order)
		engine.limitTradeOperate.buyLimitQueue.Unlock()

	} else if order.Direction == model.DirectionSellInt {
		engine.limitTradeOperate.sellLimitQueue.Lock()
		engine.limitTradeOperate.JoinInSellLimitQueue(order)
		// 加入到卖盘
		engine.sellTradePlate.Add(order)
		engine.limitTradeOperate.sellLimitQueue.Unlock()
	}
}

// 订单完成后，需要更新订单列表
func (engine *CoinTradeEngine) sendCompleteOrder(order *model.ExchangeOrder) {
	if order.Status != model.OrderStatus_Completed {
		return
	}

	marshal, _ := json.Marshal(order)

	data := kafka.KafkaData{
		Topic: topicExchangeOrderComplete,
		Key:   []byte(order.Symbol),
		Data:  marshal,
	}
	for {
		err := engine.kCli.SendSync(data)
		if err == nil {
			break
		}
		logx.Error(err)
		time.Sleep(time.Millisecond * 250)
		continue
	}
}

func NewCoinTrade(symbol string, cli *kafka.KafkaClient, db *zerodb.ZeroDB) *CoinTradeEngine {
	engine := &CoinTradeEngine{
		symbol: symbol,
		kCli:   cli,
	}
	engine.init(db)
	return engine
}
