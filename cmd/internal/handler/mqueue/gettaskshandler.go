package mqueue

import (
	"go-zero-base/utils/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mqueue/cmd/internal/logic/mqueue"
	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"
)

func GetTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TasksCollectionReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParseParamErrResponse(r, w, err)
			return
		}

		if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
			response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
			return
		}

		l := mqueue.NewGetTasksLogic(r.Context(), svcCtx)
		resp, err := l.GetTasks(&req)
		response.Response(r, w, resp, err)
	}
}
