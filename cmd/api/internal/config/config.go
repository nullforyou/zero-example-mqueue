package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
	Mysql struct {
		DataSource  string
		TablePrefix string
	}
}
