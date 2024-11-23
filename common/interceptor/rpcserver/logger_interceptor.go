package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-common/zerr"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*zerr.CodeError); ok { // 自定义错误类型
			logx.WithContext(ctx).Errorf("RPC-SERVICE-ERR | %+v", err)
			// 转成grpc err
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("RPC-SERVICE-ERR | %+v", err)
		}
	}
	return resp, err
}
