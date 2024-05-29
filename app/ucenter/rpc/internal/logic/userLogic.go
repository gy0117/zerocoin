package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/user"
	"time"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/svc"
	"ucenter-rpc/internal/verify"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	captchaVerify *verify.MachineVerify
	userDomain    *domain.UserDomain
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		captchaVerify: verify.NewMachineVerify(),
		userDomain:    domain.NewUserDomain(svcCtx.DB),
	}
}

func (l *UserLogic) FindUserById(in *user.UserRequest) (*user.UserInfoResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	defer cancel()
	result, err := l.userDomain.FindUserById(ctx, in.GetUserId())
	if err != nil {
		return nil, err
	}

	resp := &user.UserInfoResp{}

	if err = copier.Copy(resp, result); err != nil {
		return nil, err
	}

	return resp, nil
}
