package logic

import (
	"context"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"time"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) AddOrder(req *types.ExchangeReq) (string, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)

	// rpc调用
	orderResp, err := l.svcCtx.OrderRpc.AddOrder(ctx, &order.OrderReq{
		Symbol:      req.Symbol,
		Price:       req.Price,
		Amount:      req.Amount,
		Direction:   req.Direction,
		Type:        req.Type,
		UseDiscount: int32(req.UseDiscount),
		UserId:      userId,
	})
	if err != nil {
		return "", errors.Wrapf(err, "exchange-api AddOrder, uid: %d, req: %+v", userId, req)
	}
	return orderResp.OrderId, nil
}
