package handler

import (
	"github.com/zeromicro/go-zero/rest"
	"mqueue/cmd/api/internal/svc"
	"net/http"
)

func RegisterStaticFileHandler(server *rest.Server, serverCtx *svc.ServiceContext) {

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/favourite.ico",
				Handler: fileHandler("./www/favourite.ico"),
			}, {
				Method:  http.MethodGet,
				Path:    "/",
				Handler: fileHandler("./www/index.html"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/modify.html",
				Handler: fileHandler("./www/web/modify.html"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/js/jquery-1.12.4.min.js",
				Handler: fileHandler("./www/web/js/jquery-1.12.4.min.js"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/js/tool.js",
				Handler: fileHandler("./www/web/js/tool.js"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/layui/layui.js",
				Handler: fileHandler("./www/web/layui/layui.js"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/layui/css/layui.css",
				Handler: fileHandler("./www/web/layui/css/layui.css"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/layui/font/iconfont.woff2",
				Handler: fileHandler("./www/web/layui/font/iconfont.woff2"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/layui/font/iconfont.woff",
				Handler: fileHandler("./www/web/layui/font/iconfont.woff"),
			}, {
				Method:  http.MethodGet,
				Path:    "/web/layui/font/iconfont.ttf",
				Handler: fileHandler("./www/web/layui/font/iconfont.ttf"),
			},
		})
}

// 处理函数,传入文件地址
func fileHandler(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath)
	}
}
