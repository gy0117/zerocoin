package logic

import (
	"context"
	"grpc-common/market/types/rate"
	"market-rpc/internal/domain"
	"market-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeRateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeRateDomain *domain.ExchangeRateDomain
}

func NewExchangeRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeRateLogic {
	return &ExchangeRateLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		exchangeRateDomain: domain.NewExchangeRateDomain(svcCtx.Cache),
	}
}

func (l *ExchangeRateLogic) UsdRate(in *rate.RateRequest) (*rate.RateResponse, error) {
	usdRate := l.exchangeRateDomain.UsdRate(in.Unit)
	return &rate.RateResponse{
		Rate: usdRate,
	}, nil
}
