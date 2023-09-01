package svc

import (
	"github.com/hibiken/asynq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"mqueue/cmd/rpc/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	Scheduler *asynq.PeriodicTaskManager
	DbEngine  *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.Mysql.TablePrefix + "_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,                      // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		Scheduler: newPeriodicTaskManager(c, db),
		DbEngine:  db,
	}
}
