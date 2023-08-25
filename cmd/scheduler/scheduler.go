package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/scheduler/internal/config"
	"mqueue/cmd/scheduler/internal/svc"
	"os"
)

var configFile = flag.String("f", "etc/mqueue-scheduler-dev.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config
	logx.DisableStat()
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcContext := svc.NewServiceContext(c)
	//ctx := context.Background()
	//scheduler := logic.NewCronScheduler(ctx, svcContext.Scheduler)
	//scheduler.Register()

	if err := svcContext.Scheduler.Run(); err != nil {
		logx.Errorf("mqueue-scheduler运行错误 err:%+v", err)
		os.Exit(1)
	}
}
