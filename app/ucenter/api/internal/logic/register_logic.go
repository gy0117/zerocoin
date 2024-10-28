package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"grpc-common/ucenter/types/register"
	"time"

	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.Request) (*types.RegisterResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*30)
	defer cancel()

	// api和rpc模块中的Req参数是类似的
	registerReq := &register.RegisterReq{}
	if err := copier.Copy(registerReq, req); err != nil {
		return nil, err
	}

	_, err := l.svcCtx.UCRegisterRpc.RegisterByPhone(ctx, registerReq)
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	return nil, nil
}

func (l *RegisterLogic) SendCode(req *types.CodeReq) (resp *types.CodeResp, err error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*30)
	defer cancel()
	data, err := l.svcCtx.UCRegisterRpc.SendCode(ctx, &register.CodeReq{
		Country: req.Country,
		Phone:   req.Phone,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	return &types.CodeResp{
		SmsCode: data.GetSmsCode(),
	}, nil
}
