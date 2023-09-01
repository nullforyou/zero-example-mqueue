package svc

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"mqueue/cmd/dao/query"
	"mqueue/cmd/rpc/internal/config"
	"time"
)

type StorageConfigProvider struct {
	DbEngine *gorm.DB
}

func getLoadLocation() *time.Location {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return location
}

func getRedisClientOpt(c config.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
	}
}

/** newScheduler 定时任务 */
func newScheduler(c config.Config) *asynq.Scheduler {
	return asynq.NewScheduler(getRedisClientOpt(c), getSchedulerOpts())
}

/** newPeriodicTaskManager 动态定时任务 */
func newPeriodicTaskManager(c config.Config, db *gorm.DB) *asynq.PeriodicTaskManager {
	provider := &StorageConfigProvider{DbEngine: db}
	periodicTaskManager, _ := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               getRedisClientOpt(c),
			SchedulerOpts:              getSchedulerOpts(),
			PeriodicTaskConfigProvider: provider,
			SyncInterval:               60 * time.Second, //同步频率
		})
	return periodicTaskManager
}

func getSchedulerOpts() *asynq.SchedulerOpts {
	return &asynq.SchedulerOpts{
		Location: getLoadLocation(),
		PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
			if err != nil {
				logx.Errorf("定时任务排队后功能错误 task: %+v;err: %+v", info, err)
			}
		},
	}
}

func (s StorageConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	query.SetDefault(s.DbEngine)
	schedulerDao := query.Scheduler
	schedulers, _ := schedulerDao.WithContext(context.Background()).Find()

	var configs []*asynq.PeriodicTaskConfig

	limitTime, _ := time.ParseDuration("10m") //生命周期为10分钟
	deadline := time.Now().Add(limitTime)
	for _, scheduler := range schedulers {
		payload, _ := json.Marshal(scheduler)
		task := asynq.NewTask(
			scheduler.TaskType,
			payload,
			asynq.MaxRetry(0), //0重试
			asynq.Queue("critical"),
			asynq.TaskID("PeriodicTask:"+scheduler.TaskName), //独特的任务,使用TaskId避免重复
			asynq.Deadline(deadline),                         //生命周期，超过生命周期时间内未执行，将被放弃
			asynq.Timeout(30*time.Second),                    //超时时间
			asynq.
		)
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: scheduler.CronSpec, Task: task})
	}
	return configs, nil
}
