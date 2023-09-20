package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/config"
	"mqueue/cmd/scheduler"
)

var schedulerConfigFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*schedulerConfigFile, &c, conf.UseEnv())

	server := scheduler.NewPeriodicTaskManager(c, scheduler.DbEngine(c))

	//周期任务运行
	if err := server.Run(); err != nil {
		logx.Errorf("scheduler运行错误 err:%+v", err)
	}
}
