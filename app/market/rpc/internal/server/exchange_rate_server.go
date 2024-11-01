// Code generated by goctl. DO NOT EDIT.
// Source: register.proto

package server

import (
	"context"
	"grpc-common/market/types/rate"
	"market-rpc/internal/logic"
	"market-rpc/internal/svc"
)

type RateServer struct {
	svcCtx *svc.ServiceContext
	rate.UnimplementedExchangeRateServer
}

func NewRateServer(svcCtx *svc.ServiceContext) *RateServer {
	return &RateServer{
		svcCtx: svcCtx,
	}
}

func (s *RateServer) UsdRate(ctx context.Context, in *rate.RateRequest) (*rate.RateResponse, error) {
	l := logic.NewExchangeRateLogic(ctx, s.svcCtx)
	return l.UsdRate(in)
}

