package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
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
	MQueueRpc zrpc.RpcClientConf
}