package mqueue

import (
	"go-zero-base/utils/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mqueue/cmd/internal/logic/mqueue"
	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"
)

func ExecuteTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExecuteTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParseParamErrResponse(r, w, err)
			return
		}

		if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
			response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
			return
		}

		l := mqueue.NewExecuteTaskLogic(r.Context(), svcCtx)
		resp, err := l.ExecuteTask(&req)
		response.Response(r, w, resp, err)
	}
}
