package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"job-center/internal/config"
	"job-center/internal/svc"
	"job-center/internal/task"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()

	logx.MustSetup(logx.LogConf{
		Encoding: "plain",
		Stat:     false,
	})

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	t := task.NewTask(ctx)
	// 优雅退出
	go func() {
		exit := make(chan os.Signal)
		// 中断信号、终止信号
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			logx.Info("job-center中断执行，开始清理资源")
			t.Stop()
			ctx.MongoClient.DisConnect()
		}
	}()
	t.Run()
	t.StartBlocking()
}
