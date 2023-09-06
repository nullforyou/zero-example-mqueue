package handler

import (
	"github.com/zeromicro/go-zero/rest"
	"mqueue/cmd/internal/svc"
	"net/http"
	"os"
)

func getFileNames(patern string, dirPath string) []string {
	var fileNames []string
	readDir, _ := os.ReadDir(dirPath)
	for _, readFile := range readDir {
		if readFile.IsDir() {
			for _, file := range getFileNames(patern+"/"+readFile.Name(), dirPath+"/"+readFile.Name()) {
				fileNames = append(fileNames, file)
			}
		} else {
			fileNames = append(fileNames, patern+"/"+readFile.Name())
		}
	}
	return fileNames
}

func RegisterStaticFileHandler(server *rest.Server, serverCtx *svc.ServiceContext) {
	patern := "web"
	dirPath := "./www"
	fileNames := getFileNames(patern, dirPath+"/"+patern)
	for _, fileName := range fileNames {
		server.AddRoute(rest.Route{
			Method:  http.MethodGet,
			Path:    "/" + fileName,
			Handler: fileHandler("./www/" + fileName),
		})
	}

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
			},
		})
}

// 处理函数,传入文件地址
func fileHandler(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath)
	}
}
