package logic

import (
	"context"
	"errors"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/svc"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/user"
	"grpc-common/ucenter/types/wallet"
	"time"
	"zero-common/zerodb"
	"zero-common/zerodb/tran"
)

const (
	topicAddExchangeOrder = "add_exchange_order"
)

type OrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	orderDomain *domain.OrderDomain
	transaction tran.Transaction
	kafkaDomain *domain.KafkaDomain
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	orderDomain := domain.NewOrderDomain(svcCtx.DB)
	//kDomain := domain.NewKafkaDomain(svcCtx.KCli)
	//go kDomain.WaitAddOrderResult(orderDomain)

	return &OrderLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		orderDomain: orderDomain,
		transaction: tran.NewTransaction(svcCtx.DB.Conn),
		kafkaDomain: domain.NewKafkaDomain(svcCtx.KCli, orderDomain),
	}
}

func (l *OrderLogic) GetHistoryOrder(in *order.OrderReq) (*order.OrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	exchangeOrders, total, err := l.orderDomain.GetHistoryOrder(ctx, in.GetSymbol(), in.GetUserId(), in.GetPage(), in.GetPageSize())
	if err != nil {
		return nil, err
	}

	//var list []*order.ExchangeOrder
	//if err = copier.Copy(&list, exchangeOrders); err != nil {
	//	return nil, err
	//}

	list := make([]*order.ExchangeOrder, len(exchangeOrders))
	for i, v := range exchangeOrders {
		list[i] = &order.ExchangeOrder{
			Id:            v.Id,
			OrderId:       v.OrderId,
			Amount:        v.Amount,
			BaseSymbol:    v.BaseSymbol,
			CanceledTime:  v.CanceledTime,
			CoinSymbol:    v.CoinSymbol,
			CompletedTime: v.CompletedTime,
			Direction:     model.DirectionMap.Value(v.Direction),
			UserId:        v.UserId,
			Price:         v.Price,
			Status:        int32(v.Status),
			Symbol:        v.Symbol,
			Time:          v.Time,
			TradedAmount:  v.TradedAmount,
			Turnover:      v.Turnover,
			Type:          model.TypeMap.Value(v.Type),
			UseDiscount:   v.UseDiscount,
		}
	}

	resp := &order.OrderResp{
		List:  list,
		Total: total,
	}
	return resp, nil
}

func (l *OrderLogic) GetCurrentOrder(in *order.OrderReq) (*order.OrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	exchangeOrders, total, err := l.orderDomain.GetCurrentOrder(ctx, in.GetSymbol(), in.GetUserId(), in.GetPage(), in.GetPageSize())
	if err != nil {
		return nil, err
	}

	//var list []*order.ExchangeOrder
	//if err = copier.Copy(&list, exchangeOrders); err != nil {
	//	return nil, err
	//}

	list := make([]*order.ExchangeOrder, len(exchangeOrders))
	for i, v := range exchangeOrders {
		list[i] = &order.ExchangeOrder{
			Id:            v.Id,
			OrderId:       v.OrderId,
			Amount:        v.Amount,
			BaseSymbol:    v.BaseSymbol,
			CanceledTime:  v.CanceledTime,
			CoinSymbol:    v.CoinSymbol,
			CompletedTime: v.CompletedTime,
			Direction:     model.DirectionMap.Value(v.Direction),
			UserId:        v.UserId,
			Price:         v.Price,
			Status:        int32(v.Status),
			Symbol:        v.Symbol,
			Time:          v.Time,
			TradedAmount:  v.TradedAmount,
			Turnover:      v.Turnover,
			Type:          model.TypeMap.Value(v.Type),
			UseDiscount:   v.UseDiscount,
		}
	}

	resp := &order.OrderResp{
		List:  list,
		Total: total,
	}
	return resp, nil
}

// 发布委托：就是创建订单，一旦发布，就需要冻结 钱和手续费

// AddOrder 逻辑
// 1. 判断参数是否合法，判断钱是否足够
// 2. 创建订单，得到订单id和需要冻结的钱
// 3. 发送消息到kafka
// 4. ucenter接收消息，处理钱包冻结
// 5. 失败则调用订单服务，取消订单
func (l *OrderLogic) AddOrder(req *order.OrderReq) (*order.AddOrderResp, error) {
	// 1. 参数合法性检查
	if req.GetType() == model.LimitPrice && req.GetPrice() <= 0 {
		return nil, errors.New("in limit order mode, the price cannot be less than or equal to 0")
	}
	if req.GetAmount() <= 0 {
		return nil, errors.New("the purchase quantity cannot be less than or equal to 0")
	}

	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// 2. 查询用户的状态，是否禁止买卖（查询user表）
	userResp, err := l.svcCtx.UserRpc.FindUserById(ctx, &user.UserRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	// 交易状态，0 禁止交易
	if userResp.TransactionStatus == 0 {
		return nil, errors.New("current user must be forbidden trade")
	}

	// 3. 查表exchange_coin，对应的symbol是否可以交易
	exchangeCoinResp, err := l.svcCtx.MarketRpc.FindSymbolInfo(ctx, &market.MarketRequest{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, err
	}
	// enable 状态，1：启用，2：禁止   exchangeable是否可交易，1：可交易
	if exchangeCoinResp.Enable != 1 && exchangeCoinResp.Exchangeable != 1 {
		return nil, errors.New(fmt.Sprintf("%s not tradable", req.Symbol))
	}

	// 4. 查询待买入卖出的币是否支持（coin表中是否存在这个币）
	coinUnit := exchangeCoinResp.GetBaseSymbol() // 基准币，例如：USDT
	if req.Direction == model.DirectionSell {
		coinUnit = exchangeCoinResp.GetCoinSymbol() // 交易币，例如：BTC
	}
	_, err = l.svcCtx.MarketRpc.FindCoinInfo(ctx, &market.MarketRequest{
		Unit: coinUnit,
	})
	if err != nil {
		return nil, err
	}

	// 5. 查询用户钱包（基准币和交易币都查询）
	baseSymbolWalletResp, err := l.svcCtx.WalletRpc.FindWalletBySymbol(ctx, &wallet.WalletReq{
		CoinName: exchangeCoinResp.GetBaseSymbol(),
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, err
	}
	coinSymbolWalletResp, err := l.svcCtx.WalletRpc.FindWalletBySymbol(ctx, &wallet.WalletReq{
		CoinName: exchangeCoinResp.GetCoinSymbol(),
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, err
	}
	// 钱包是否锁定 0 否 1 是
	if baseSymbolWalletResp.IsLock == 1 || coinSymbolWalletResp.IsLock == 1 {
		return nil, errors.New(fmt.Sprintf("%d wallet is locked", req.UserId))
	}

	// 6. 限制委托数量
	// 查询当前用户正在交易的数量不能超过最大交易数量
	count, err := l.orderDomain.FindCurrentTradingCount(ctx, req.UserId, req.Symbol, req.Direction)
	if err != nil {
		return nil, err
	}
	if exchangeCoinResp.GetMaxTradingOrder() > 0 && count > exchangeCoinResp.GetMaxTradingOrder() {
		return nil, errors.New(fmt.Sprintf("exceeds the maximum order quantity %d", exchangeCoinResp.GetMaxTradingOrder()))
	}

	// 7. 保存订单到数据库，发送消息到kafka，【ucenter钱包服务】 接收到消息，进行资金的冻结
	// 如果消息发送失败，则整体回滚

	// 生成订单
	newOrder := model.NewOrder()
	newOrder.UserId = req.UserId
	newOrder.Symbol = req.Symbol
	newOrder.Type = model.TransferType(req.Type)
	newOrder.Direction = model.TransferDirection(req.Direction)
	newOrder.BaseSymbol = exchangeCoinResp.GetBaseSymbol()
	newOrder.CoinSymbol = exchangeCoinResp.GetCoinSymbol()

	// 市价：会以当前市场价格快速成交，这个价格是变动的，表中不会记录
	// 限价：用户自己输入的，用户自己的期望价格
	if req.Type == model.MarketPrice {
		newOrder.Price = 0
	} else {
		newOrder.Price = req.Price
	}
	newOrder.UseDiscount = "0"
	newOrder.Amount = req.Amount

	err = l.transaction.Action(func(conn zerodb.DbConn) error {
		// 1. 提交订单
		money, err := l.orderDomain.AddOrder(ctx, conn, newOrder, baseSymbolWalletResp, coinSymbolWalletResp)
		if err != nil {
			return errors.New("order submission failed")
		}
		// 2. 发送消息到kafka，冻结钱包服务的消息
		err = l.kafkaDomain.SendOrderAdd(topicAddExchangeOrder, req.UserId, newOrder.OrderId, money, req.Symbol, newOrder.Direction, exchangeCoinResp.GetBaseSymbol(), exchangeCoinResp.GetCoinSymbol())
		if err != nil {
			return errors.New("the order message failed to send")
		}
		return nil
	})

	return &order.AddOrderResp{
		OrderId: newOrder.OrderId,
	}, nil
}

func (l *OrderLogic) FindByOrderId(req *order.OrderReq) (*order.ExchangeOrder, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()
	exchangeOrder, err := l.orderDomain.FindByOrderId(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("exchange | FindByOrderId | exchangeOrder: %+v\n", exchangeOrder)
	resp := &order.ExchangeOrder{}
	if err = copier.Copy(resp, exchangeOrder); err != nil {
		return nil, err
	}

	fmt.Printf("exchange | FindByOrderId | resp: %+v\n", resp)
	return resp, nil
}

func (l *OrderLogic) CancelOrder(req *order.OrderReq) (*order.CancelOrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	err := l.orderDomain.CancelOrder(ctx, req)
	return nil, err
}
