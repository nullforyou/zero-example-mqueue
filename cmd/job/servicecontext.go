package job

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"mqueue/cmd/config"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.Redis)
	return &ServiceContext{
		Config: c,
		Redis:  rds,
	}
}
