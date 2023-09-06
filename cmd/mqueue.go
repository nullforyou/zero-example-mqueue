package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-base/utils/response"
	"go-zero-base/utils/xerr"
	"net/http"

	"mqueue/cmd/config"
	"mqueue/cmd/internal/handler"
	"mqueue/cmd/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	flag.Parse()
	logx.DisableStat()
	var c config.Config
	conf.MustLoad(*config.GetConfigFile(), &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		//JWT验证失败自定义处理
		response.Response(r, w, nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorTokenExpire)))
	}))

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//注册目录处理器
	handler.RegisterStaticFileHandler(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
