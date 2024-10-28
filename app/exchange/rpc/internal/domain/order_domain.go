package domain

import (
	"context"
	"errors"
	"exchange-rpc/internal/dao"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/repo"
	"grpc-common/exchange/types/order"
	"grpc-common/ucenter/uclient"
	"time"
	"zero-common/operate"
	"zero-common/tools"
	"zero-common/zerodb"
)

type OrderDomain struct {
	orderRepo repo.OrderRepo
}

func NewOrderDomain(db *zerodb.ZeroDB) *OrderDomain {
	return &OrderDomain{
		orderRepo: dao.NewOrderDao(db),
	}
}

func (d *OrderDomain) GetHistoryOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) ([]*model.ExchangeOrder, int64, error) {
	list, total, err := d.orderRepo.GetHistoryOrder(ctx, symbol, memId, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		return nil, 0, errors.New("data not found")
	}
	return list, total, nil
}

func (d *OrderDomain) GetCurrentOrder(ctx context.Context, symbol string, memId int64, pageNo int64, pageSize int64) ([]*model.ExchangeOrder, int64, error) {
	list, total, err := d.orderRepo.GetCurrentOrder(ctx, symbol, memId, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		return nil, 0, errors.New("data not found")
	}
	return list, total, nil
}

func (d *OrderDomain) FindCurrentTradingCount(ctx context.Context, memId int64, symbol string, direction string) (int64, error) {
	return d.orderRepo.FindCurrentTradingCount(ctx, memId, symbol, model.TransferDirection(direction))
}

func (d *OrderDomain) AddOrder(
	ctx context.Context,
	conn zerodb.DbConn,
	order *model.ExchangeOrder,
	baseWallet *uclient.WalletResp,
	coinWallet *uclient.WalletResp) (float64, error) {

	order.Status = model.OrderStatus_StatusInit
	order.Time = time.Now().UnixMilli()
	order.OrderId = tools.GenerateOrderId("eo")

	//交易的时候  coin.Fee 费率 手续费 我们做的时候 先不考虑手续费
	//买 花USDT 市价 price 0 冻结的直接就是amount  卖 BTC
	var money float64

	if order.Direction == model.DirectionBuyInt {
		// 买入
		if order.Type == model.MarketPriceInt { // 市价
			money = order.Amount
		} else {
			money = operate.FloorFloat(order.Price*order.Amount, 8)
		}
		if baseWallet.Balance < money {
			return 0, errors.New("insufficient balance")
		}
	} else {
		money = order.Amount
		if coinWallet.Balance < money {
			return 0, errors.New("insufficient balance")
		}
	}

	err := d.orderRepo.Save(ctx, conn, order)
	return money, err
}

func (d *OrderDomain) FindByOrderId(ctx context.Context, req *order.OrderReq) (*model.ExchangeOrder, error) {
	eo, err := d.orderRepo.FindByOrderId(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	//if eo == nil {
	//	return nil, errors.New("订单不存在")
	//}
	return eo, nil
}

func (d *OrderDomain) CancelOrder(ctx context.Context, req *order.OrderReq) error {
	return d.orderRepo.CancelOrder(ctx, req.OrderId)
}

func (d *OrderDomain) UpdateOrderStatus(ctx context.Context, orderId string, status int) error {
	return d.orderRepo.UpdateOrderStatus(ctx, orderId, status)
}

// FindTradingOrderBySymbol 查询 正在交易的且是symbol交易对的 的订单
func (d *OrderDomain) FindTradingOrderBySymbol(ctx context.Context, symbol string) ([]*model.ExchangeOrder, error) {
	return d.orderRepo.FindTradingOrderBySymbol(ctx, symbol)
}

func (d *OrderDomain) UpdateOrderComplete(ctx context.Context, order *model.ExchangeOrder) error {
	return d.orderRepo.UpdateOrderComplete(ctx, order.OrderId, order.TradedAmount, order.Turnover)
}
