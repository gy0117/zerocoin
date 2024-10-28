package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/chain"
	"market-api/internal/config"
	"market-api/internal/handler"
	"market-api/internal/svc"
	"market-api/internal/ws"
	"net/http"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.LogConfig)

	wsServer := ws.NewWebSocketServer("/socket.io")

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithChain(chain.New(wsServer.ServerHandler)), // 中间件
		rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,X-Auth-Token,x-auth-token")
		}, nil, "http://localhost:8080"))

	defer server.Stop()

	ctx := svc.NewServiceContext(c, wsServer)
	handler.RegisterHandlers(server, ctx)

	// 多个server一起启动
	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(wsServer)

	logx.Infof("Starting api server at %s:%d...", c.Host, c.Port)

	group.Start()
}
