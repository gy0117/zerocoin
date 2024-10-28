package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/user"
	"time"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	"zero-common/tools"
)

type ApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLogic {
	return &ApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLogic) CheckSecuritySetting() (*types.ApproveResp, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()

	userId := ctx.Value("userId").(int64)
	memberResp, err := l.svcCtx.UserRpc.FindUserById(ctx, &user.UserRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "uid: %d", userId)
	}

	resp := &types.ApproveResp{}
	resp.Username = memberResp.Username
	resp.CreateTime = tools.ToTimeString(memberResp.RegistrationTime)
	resp.Id = memberResp.Id

	if memberResp.Email != "" {
		resp.EmailVerified = "true"
		resp.Email = memberResp.Email
	} else {
		resp.EmailVerified = "false"
	}
	if memberResp.JyPassword != "" {
		resp.FundsVerified = "true"
	} else {
		resp.FundsVerified = "false"
	}
	resp.LoginVerified = "true"
	if memberResp.MobilePhone != "" {
		resp.PhoneVerified = "true"
		resp.MobilePhone = memberResp.MobilePhone
	} else {
		resp.PhoneVerified = "false"
	}
	if memberResp.RealName != "" {
		resp.RealVerified = "true"
		resp.RealName = memberResp.RealName
	} else {
		resp.RealVerified = "false"
	}
	resp.IdCard = memberResp.IdNumber
	if memberResp.IdNumber != "" {
		resp.IdCard = memberResp.IdNumber[:2] + "********"
	}
	//0 未认证 1 审核中 2 已认证
	if memberResp.RealNameStatus == 1 {
		resp.RealAuditing = "true"
	} else {
		resp.RealAuditing = "false"
	}
	resp.Avatar = memberResp.Avatar
	if memberResp.Bank == "" && memberResp.AliNo == "" && memberResp.Wechat == "" {
		resp.AccountVerified = "false"
	} else {
		resp.AccountVerified = "true"
	}
	return resp, nil
}
