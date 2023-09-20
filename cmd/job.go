package main

import (
	"context"
	"flag"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/business"
	"mqueue/cmd/config"
	"mqueue/cmd/job"
	"os"
	"strings"
)

var jobConfigFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	//logx.DisableStat()
	var c config.Config
	conf.MustLoad(*jobConfigFile, &c, conf.UseEnv())

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := job.NewServiceContext(c)
	mux := asynq.NewServeMux()
	mux.Handle(business.PlatformHttp, job.NewPlatformHttpHandler(ctx))

	server := newAsynqServer(c)
	if err := server.Run(mux); err != nil {
		logx.Errorf("== >> [%s] 运行作业错误 err:%+v", c.Name, err)
		os.Exit(1)
	}
}

func newAsynqServer(c config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass},
		asynq.Config{
			IsFailure: func(err error) bool {
				if strings.EqualFold(err.Error(), "context deadline exceeded") {
					//超过任务有效期，不是一个错误，返回 false
					return false
				} else {
					logx.Errorf("asynq服务器执行任务失败 err : %+v \n", err)
					return true
				}
			},
			Concurrency: 20, //最大并发进程任务数
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: false, //关键队列中的任务总是首先被处理。如果关键队列为空，则处理默认队列。如果关键队列和默认队列都为空，则处理低队列。
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				//if strings.EqualFold(err.Error(), "context deadline exceeded") {
				logx.Infof("任务超过有效期,payload:%s;", task.Payload())
				//}
			}),
		},
	)
}
