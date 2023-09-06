package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"mqueue/cmd/config"
	"mqueue/cmd/dao/query"
	"time"
)

type StorageConfigProvider struct {
	DbEngine *gorm.DB
}

func DbEngine(c config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.Mysql.TablePrefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,                // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	return db
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

// NewScheduler 定时任务
func NewScheduler(c config.Config) *asynq.Scheduler {
	return asynq.NewScheduler(getRedisClientOpt(c), getSchedulerOpts())
}

// NewPeriodicTaskManager 动态定时任务
func NewPeriodicTaskManager(c config.Config, db *gorm.DB) *asynq.PeriodicTaskManager {
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
				if !errors.Is(err, asynq.ErrTaskIDConflict) { //表示非队列Id重复错误
					logx.Errorf("定时任务排队后功能错误 task: %+v;err: %+v", info, err)
				}
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
		)
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: scheduler.CronSpec, Task: task})
	}
	return configs, nil
}
