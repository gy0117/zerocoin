package logic

import (
	"common/tools"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/login"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*30)
	defer cancel()

	// 这里的参数需要转换，api层和rpc层的对象不要用同一个
	loginReq := &login.LoginReq{}
	_ = copier.Copy(loginReq, req)

	resp, err := l.svcCtx.UCLoginRpc.Login(ctx, loginReq)
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	result := &types.LoginResp{}
	_ = copier.Copy(result, resp)
	return result, nil
}

func (l *LoginLogic) CheckLogin(token string, f func() bool) (bool, error) {
	// jwt 解析token
	_, err := tools.ParseToken(token, l.svcCtx.Config.Jwt.AccessSecret)

	if err != nil {
		return false, errors.Wrapf(err, "token: %s", token)
	}

	if f() {
		return false, nil
	}
	return true, nil
}

func (l *LoginLogic) RefreshToken(refreshToken string) (string, error) {
	userId, err := tools.ParseToken(refreshToken, l.svcCtx.Config.Jwt.AccessSecret)

	if err != nil {
		return "", errors.Wrapf(err, "parse refresh-token: %s", refreshToken)
	}

	newAccessToken, err := tools.GenerateAccessToken(userId)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate access token: %s", refreshToken)
	}
	return newAccessToken, nil
}
