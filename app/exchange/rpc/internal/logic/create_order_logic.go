package logic

import (
	"context"
	"database/sql"
	"exchange-rpc/internal/config"
	"exchange-rpc/internal/domain"
	domain2 "exchange-rpc/internal/domain/v2"
	"exchange-rpc/internal/model"
	"exchange-rpc/internal/svc"
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-common/exchange/types/order"
	"zero-common/dtmutil"
	"zero-common/operate"
	"zero-common/zerodb/tran"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	orderDomain   *domain.OrderDomain
	orderDomainV2 *domain2.OrderDomainV2
	transaction   tran.Transaction
	kafkaDomain   *domain.KafkaDomain
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	orderDomain := domain.NewOrderDomain(svcCtx.DB)

	return &CreateOrderLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		orderDomain:   orderDomain,
		orderDomainV2: domain2.NewOrderDomainV2(),
		transaction:   tran.NewTransaction(svcCtx.DB.Conn),
		kafkaDomain:   domain.NewKafkaDomain(svcCtx.KCli, orderDomain),
	}
}

func dbGet(c config.Config) *dtmutil.DB {
	var dbConf = dtmcli.DBConf{
		Driver:   c.Mysql.Driver,
		Host:     c.Mysql.Host,
		Port:     c.Mysql.Port,
		User:     c.Mysql.User,
		Password: c.Mysql.Password,
		Db:       c.Mysql.Db,
	}
	return dtmutil.DbGet(dbConf)
}

// CreateOrder 只和创建订单有关的逻辑，前期的check不在函数职能以内
func (l *CreateOrderLogic) CreateOrder(createOrderReq *order.CreateOrderRequest) (*order.AddOrderResp, error) {
	logx.Infof("saga -> 创建订单, createOrderReq: %+v", createOrderReq)

	var money float64
	item := createOrderReq.Item
	if item.Direction == model.DirectionBuy {
		if item.Type == model.MarketPrice { // 市价买入
			money = item.Amount
		} else {
			// 限价买入
			money = operate.FloorFloat(item.Price*item.Amount, 8)
		}
		// 钱不够
		if createOrderReq.GetBaseBalance() < money {
			return nil, errors.New("not enough balance when create order")
		}
	} else {
		// 卖出，卖出btc，钱就是个数
		money = item.Amount
		if createOrderReq.GetCoinBalance() < money {
			return nil, errors.New("not enough balance when create order")
		}
	}

	// grpc order -> sql order
	newOrder := model.NewExchangeOrder(createOrderReq)

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, errors.Wrap(status.Error(codes.Aborted, err.Error()), "create_order_logic create barrier failed")
	}
	tx := dbGet(l.svcCtx.Config).DB.Begin()
	//sourceTx := tx.Statement.ConnPool.(*gorm.PreparedStmtTX).Tx.(*sql.Tx)
	sourceTx := tx.Statement.ConnPool.(*sql.Tx)

	if err := barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
		err := l.orderDomainV2.CreateOrder(l.ctx, tx, newOrder)
		if err != nil {
			return fmt.Errorf("create order failed, newOrder: %+v, err: %v", newOrder, err)
		}
		return nil
	}); err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.AddOrderResp{
		OrderId: newOrder.OrderId,
	}, nil
}

// TODO
func (l *CreateOrderLogic) CreateOrderRevert(req *order.CreateOrderRequest) (*order.AddOrderResp, error) {
	logx.Error("saga -> 创建订单回滚")
	return &order.AddOrderResp{}, nil
}
