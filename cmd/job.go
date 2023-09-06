package main

import (
	"context"
	"flag"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/business"
	"mqueue/cmd/config"
	"os"
)

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*config.GetConfigFile(), &c, conf.UseEnv())

	// log、prometheus、trace、metricsUrl
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	mux := asynq.NewServeMux()
	mux.HandleFunc(business.PlatformHttp, PlatformHttpHandler)

	server := newAsynqServer(c)
	if err := server.Run(mux); err != nil {
		logx.Errorf("== >> [%s] 运行作业错误 err:%+v", c.Name, err)
		os.Exit(1)
	}
}

func PlatformHttpHandler(ctx context.Context, t *asynq.Task) error {
	logx.Infof("执行【%s】 payload:%s", business.PlatformHttp, t.Payload())
	return nil
	//return xerr.NewBusinessError(xerr.SetCode(xerr.ErrorServerCommon), xerr.SetMsg("测试Job错误"))
}

func newAsynqServer(c config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass},
		asynq.Config{
			IsFailure: func(err error) bool {
				logx.Errorf("asynq服务器执行任务失败 err : %+v \n", err)
				return true
			},
			Concurrency: 20, //最大并发进程任务数
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: false, //关键队列中的任务总是首先被处理。如果关键队列为空，则处理默认队列。如果关键队列和默认队列都为空，则处理低队列。
		},
	)
}
