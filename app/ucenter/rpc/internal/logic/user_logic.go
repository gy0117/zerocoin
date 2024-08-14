package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/user"
	"ucenter-rpc/internal/domain"
	"ucenter-rpc/internal/svc"
	"ucenter-rpc/internal/verify"
	"zero-common/zerr"
)

var ErrFindUser = zerr.NewCodeErr(zerr.FIND_USER_ERROR)

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
	//ctx, cancel := context.WithTimeout(l.ctx, time.Second*5)
	//defer cancel()
	result, err := l.userDomain.FindUserById(l.ctx, in.GetUserId())
	if err != nil {
		return nil, errors.Wrapf(ErrFindUser, "查找用户失败 uid: %s, err: %v", in.GetUserId(), err)
	}

	resp := &user.UserInfoResp{}

	_ = copier.Copy(resp, result)

	return resp, nil
}
