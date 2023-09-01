package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"mqueue/cmd/pb/task"
	"mqueue/cmd/rpc/internal/config"
	"mqueue/cmd/rpc/internal/server"
	"mqueue/cmd/rpc/internal/svc"
)

var configFile = flag.String("f", "etc/mqueue-rpc.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	//周期任务运行
	/*if err := ctx.Scheduler.Run(); err != nil {
		logx.Errorf("scheduler运行错误 err:%+v", err)
	}*/

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		task.RegisterTaskServer(grpcServer, server.NewTaskServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()

}
