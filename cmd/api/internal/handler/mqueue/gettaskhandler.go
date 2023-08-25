package mqueue

import (
	"go-zero-base/utils/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mqueue/cmd/api/internal/logic/mqueue"
	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"
)

func GetTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskItemReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParseParamErrResponse(r, w, err)
			return
		}

		if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
			response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
			return
		}

		l := mqueue.NewGetTaskLogic(r.Context(), svcCtx)
		resp, err := l.GetTask(&req)
		response.Response(r, w, resp, err)
	}
}
