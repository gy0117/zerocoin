package logic

import (
	"context"
	"exchange-rpc/internal/domain"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/svc"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"zero-common/operate"
	"zero-common/tools"

	// 下面这行导入gozero的dtm驱动
	//_ "github.com/dtm-labs/driver-gozero"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/user"
	"grpc-common/ucenter/types/wallet"
	"time"
	"zero-common/zerodb/tran"
	"zero-common/zerr"
)

const (
	topicAddExchangeOrder = "add_exchange_order"
)

var ErrGetHistoryOrder = zerr.NewCodeErr(zerr.EXCHANGE_GET_HISTORY_ORDER_ERROR)
var ErrGetCurrentOrder = zerr.NewCodeErr(zerr.EXCHANGE_GET_CURRENT_ORDER_ERROR)
var ErrAddOrder = zerr.NewCodeErr(zerr.EXCHANGE_ADD_ORDER_ERROR)
var ErrFindOrder = zerr.NewCodeErr(zerr.EXCHANGE_FIND_ORDER_ERROR)
var ErrCancelOrder = zerr.NewCodeErr(zerr.EXCHANGE_CANCEL_ORDER_ERROR)

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
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc in limit order mode, the price cannot be less than or equal to 0")
	}
	if req.GetAmount() <= 0 {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc the purchase amount cannot be less than or equal to 0")
	}

	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// 2. 查询用户的状态，是否禁止买卖（查询user表）
	userResp, err := l.svcCtx.UserRpc.FindUserById(ctx, &user.UserRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc findUserById, uid: %d, err: %v", req.UserId, err)
	}
	// 交易状态，0 禁止交易
	if userResp.TransactionStatus == 0 {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc current user must be forbidden trade")
	}

	// 3. 查表exchange_coin，对应的symbol是否可以交易
	exchangeCoinResp, err := l.svcCtx.MarketRpc.FindSymbolInfo(ctx, &market.MarketRequest{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc findSymbolInfo symbol: %s", req.Symbol)
	}
	// enable 状态，1：启用，2：禁止   exchangeable是否可交易，1：可交易
	if exchangeCoinResp.Enable != 1 && exchangeCoinResp.Exchangeable != 1 {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc %s not tradable", req.Symbol)
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
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc findCoinInfo, unit: %s, err: %v", coinUnit, err)
	}

	// 5. 查询用户钱包（基准币和交易币都查询）
	baseSymbolWalletResp, err := l.svcCtx.WalletRpc.FindWalletBySymbol(ctx, &wallet.WalletReq{
		CoinName: exchangeCoinResp.GetBaseSymbol(),
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder,
			"exchange-rpc findWalletBySymbol, uid: %d, coinName: %s, err: %v", req.UserId, exchangeCoinResp.GetBaseSymbol(), err)
	}
	coinSymbolWalletResp, err := l.svcCtx.WalletRpc.FindWalletBySymbol(ctx, &wallet.WalletReq{
		CoinName: exchangeCoinResp.GetCoinSymbol(),
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder,
			"exchange-rpc findWalletBySymbol, uid: %d, coinName: %s, err: %v", req.UserId, exchangeCoinResp.GetCoinSymbol(), err)
	}
	// 钱包是否锁定 0 否 1 是
	if baseSymbolWalletResp.IsLock == 1 || coinSymbolWalletResp.IsLock == 1 {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc %d wallet is locked", req.UserId)
	}

	// 6. 限制委托数量
	// 查询当前用户正在交易的数量不能超过最大交易数量
	count, err := l.orderDomain.FindCurrentTradingCount(ctx, req.UserId, req.Symbol, req.Direction)
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc findCurrentTradingCount uid: %d, symbol: %s, direction: %s", req.UserId, req.Symbol, req.Direction)
	}
	if exchangeCoinResp.GetMaxTradingOrder() > 0 && count > exchangeCoinResp.GetMaxTradingOrder() {
		return nil, errors.Wrapf(ErrAddOrder, "exchange-rpc exceeds the maximum order amount %d", exchangeCoinResp.GetMaxTradingOrder())
	}

	// 7. 保存订单到数据库，发送消息到kafka，【ucenter钱包服务】 接收到消息，进行资金的冻结
	// 如果消息发送失败，则整体回滚

	logx.Info("=============开始dtm")
	// DTM 改造
	dtmServer, err := l.svcCtx.Config.DtmConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "get dtm server failed, err: %v", err)
	}
	logx.Info("dtmServer -> ", dtmServer)
	gid := dtmgrpc.MustGenGid(dtmServer)
	logx.Info("gid -> ", gid)
	dtmSaga := dtmgrpc.NewSagaGrpc(dtmServer, gid)

	orderTarget, err := buildTarget(l.svcCtx)
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "create order failed, err: %v", err)
	}
	accountTarget, err := l.svcCtx.Config.UCenter.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(ErrAddOrder, "create order failed, err: %v", err)
	}

	var createOrderAddr = orderTarget + order.OrderService_CreateOrder_FullMethodName
	var createOrderRevertAddr = orderTarget + order.OrderService_CreateOrderRevert_FullMethodName
	var freezeUserWalletAddr = accountTarget + wallet.Wallet_FreezeUserAsset_FullMethodName
	var freezeUserWalletRevertAddr = accountTarget + wallet.Wallet_UnFreezeUserAsset_FullMethodName
	var send2PlateAddr = orderTarget + order.OrderService_SendOrder2Plate_FullMethodName
	var send2PlateRevertAddr = orderTarget + order.OrderService_SendOrder2PlateRevert_FullMethodName

	logx.Info("createOrderAddr: ", createOrderAddr)
	logx.Info("createOrderRevertAddr: ", createOrderRevertAddr)
	logx.Info("freezeUserWalletAddr: ", freezeUserWalletAddr)
	logx.Info("freezeUserWalletRevertAddr: ", freezeUserWalletRevertAddr)
	logx.Info("send2PlateAddr: ", send2PlateAddr)
	logx.Info("send2PlateRevertAddr: ", send2PlateRevertAddr)

	orderId := tools.GenerateOrderId("eo")
	item := &order.CreateOrderItem{
		UserId:      req.UserId,
		Symbol:      req.Symbol,
		Price:       req.Price,
		Amount:      req.Amount,
		Direction:   req.Direction,
		Type:        req.Type,
		UseDiscount: req.UseDiscount,
		OrderId:     orderId,
	}
	createOrderRequest := &order.CreateOrderRequest{
		Item:        item,
		BaseSymbol:  exchangeCoinResp.GetBaseSymbol(),
		BaseBalance: baseSymbolWalletResp.GetBalance(),
		CoinSymbol:  exchangeCoinResp.GetCoinSymbol(),
		CoinBalance: coinSymbolWalletResp.GetBalance(),
	}

	var money float64
	var symbol string
	if item.Direction == model.DirectionBuy {
		symbol = exchangeCoinResp.GetBaseSymbol()
		if item.Type == model.MarketPrice {
			money = item.Amount
		} else {
			money = operate.FloorFloat(item.Price*item.Amount, 8)
		}
	} else {
		symbol = exchangeCoinResp.GetCoinSymbol()
		// 卖出，卖出btc，钱就是个数
		money = item.Amount
	}

	freezeRequest := &wallet.FreezeUserAssetReq{
		Uid:    req.UserId,
		Money:  money,
		Symbol: symbol,
	}

	send2PlatRequest := &order.SendOrderRequest{
		OrderId: orderId,
	}

	saga := dtmSaga.
		Add(createOrderAddr, createOrderRevertAddr, createOrderRequest).
		Add(freezeUserWalletAddr, freezeUserWalletRevertAddr, freezeRequest).
		Add(send2PlateAddr, send2PlateRevertAddr, send2PlatRequest)
	err = saga.Submit()
	if err != nil {
		logx.Error("saga, err: ", err)
		return nil, err
	}
	return &order.AddOrderResp{}, nil
}

func buildTarget(ctx *svc.ServiceContext) (string, error) {
	etcd := ctx.Config.Etcd
	if len(etcd.Hosts) == 0 {
		return "", errors.New("build target failed")
	}
	host := etcd.Hosts[0]
	key := etcd.Key
	target := fmt.Sprintf("etcd://%s/%s", host, key)
	return target, nil
}

func (l *OrderLogic) FindByOrderId(req *order.OrderReq) (*order.ExchangeOrder, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()
	exchangeOrder, err := l.orderDomain.FindByOrderId(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(ErrFindOrder, "exchange-rpc findOrder req: %+v, err: %v", req, err)
	}
	resp := &order.ExchangeOrder{}
	_ = copier.Copy(resp, exchangeOrder)
	return resp, nil
}

func (l *OrderLogic) CancelOrder(req *order.OrderReq) (*order.CancelOrderResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	err := l.orderDomain.CancelOrder(ctx, req)
	return nil, errors.Wrapf(ErrCancelOrder, "exchange-rpc cancel order req: %+v, err: %v", req, err)
}
