// Code generated by goctl. DO NOT EDIT.
// Source: register.proto

package server

import (
	"context"
	"grpc-common/ucenter/types/login"

	"ucenter-rpc/internal/logic"
	"ucenter-rpc/internal/svc"
)

type LoginServer struct {
	svcCtx *svc.ServiceContext
	login.UnimplementedLoginServer
}

func NewLoginServer(svcCtx *svc.ServiceContext) *LoginServer {
	return &LoginServer{
		svcCtx: svcCtx,
	}
}


func (s *LoginServer) Login(ctx context.Context, in *login.LoginReq) (*login.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}