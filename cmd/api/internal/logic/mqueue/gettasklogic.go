package mqueue

import (
	"context"

	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskLogic) GetTask(req *types.TaskItemReq) (resp *types.TaskItemResp, err error) {
	// todo: add your logic here and delete this line

	return
}
