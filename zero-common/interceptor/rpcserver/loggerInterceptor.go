package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-common/zeroerr"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		// err类型
		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zeroerr.CodeError); ok {
			// 自定义错误类型
			logx.WithContext(ctx).Errorf("[RPC-SERVICE-ERR] %+v", err)
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("[RPC-SERVICE-ERR] %+v", err)
		}

	}
	return resp, err
}
