package config

import (
	"flag"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Redis redis.RedisConf
	Jwt   struct {
		AccessSecret string
		AccessExpire int64
	}
	Mysql struct {
		DataSource  string
		TablePrefix string
	}
}

func GetConfigFile() *string {
	return flag.String("f", "etc/mqueue.yaml", "the config file")
}
