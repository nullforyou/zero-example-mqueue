package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-base/custom_validate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"mqueue/cmd/api/internal/config"
	"mqueue/cmd/taskclient"
)

type ServiceContext struct {
	Config    config.Config
	DbEngine  *gorm.DB
	Validator custom_validate.Validator //验证器
	MQueueRpc taskclient.Task
}

func NewServiceContext(c config.Config) *ServiceContext {
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

	return &ServiceContext{
		Config:    c,
		DbEngine:  db,
		MQueueRpc: taskclient.NewTask(zrpc.MustNewClient(c.MQueueRpc)),
		Validator: custom_validate.InitValidator(),
	}
}
