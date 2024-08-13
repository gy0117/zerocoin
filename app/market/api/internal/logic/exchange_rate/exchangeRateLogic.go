package exchange_rate

import (
	"context"
	"github.com/pkg/errors"
	"grpc-common/market/types/rate"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeRateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeRateLogic {
	return &ExchangeRateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeRateLogic) UsdRate(req *types.RateRequest) (resp *types.RateResponse, err error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	// api和rpc模块中的Req参数是类似的
	rateReq := &rate.RateRequest{
		Unit: req.Unit,
		Ip:   req.Ip,
	}

	result, err := l.svcCtx.ExchangeRateRpc.UsdRate(ctx, rateReq)
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	return &types.RateResponse{
		Rate: result.Rate,
	}, nil
}
