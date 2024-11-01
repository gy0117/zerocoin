package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/login"
	"grpc-common/ucenter/types/user"
	"grpc-common/ucenter/types/wallet"
	"grpc-common/ucenter/types/withdraw"
	"zero-common/interceptor/rpcserver"

	"grpc-common/ucenter/types/register"
	"ucenter-rpc/internal/config"
	"ucenter-rpc/internal/server"
	"ucenter-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConfig)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		register.RegisterRegisterServer(grpcServer, server.NewRegisterServer(ctx))
		login.RegisterLoginServer(grpcServer, server.NewLoginServer(ctx))
		wallet.RegisterWalletServer(grpcServer, server.NewWalletServer(ctx))
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		withdraw.RegisterWithdrawServer(grpcServer, server.NewWithdrawServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// rpc log
	s.AddUnaryInterceptors(rpcserver.LoggerInterceptor)

	defer s.Stop()

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
