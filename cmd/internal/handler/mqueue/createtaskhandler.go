package mqueue

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-base/utils/response"
	"mqueue/cmd/internal/logic/mqueue"
	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"
	"net/http"
)

func CreateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParseParamErrResponse(r, w, err)
			return
		}

		if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
			response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
			return
		}

		l := mqueue.NewCreateTaskLogic(r.Context(), svcCtx)
		resp, err := l.CreateTask(&req)
		response.Response(r, w, resp, err)
	}
}
