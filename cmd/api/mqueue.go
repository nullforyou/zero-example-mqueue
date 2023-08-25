package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-base/utils/response"
	"go-zero-base/utils/xerr"
	"net/http"

	"mqueue/cmd/api/internal/config"
	"mqueue/cmd/api/internal/handler"
	"mqueue/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/mqueue-api.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		//JWT验证失败自定义处理
		response.Response(r, w, nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorTokenExpire)))
	}))

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
