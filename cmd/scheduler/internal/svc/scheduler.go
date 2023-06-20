package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/scheduler/internal/config"
	"time"
)

func newScheduler(c config.Config) *asynq.Scheduler {

	location,_ := time.LoadLocation("Asia/Shanghai")

	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr: c.Redis.Host,
			Password: c.Redis.Pass,
		},
		&asynq.SchedulerOpts{
			Location: location,
			PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
				if err != nil {
					logx.Errorf("定时任务排队后功能错误 task: %+v;err: %+v", info, err)
				}
			},
		})
}