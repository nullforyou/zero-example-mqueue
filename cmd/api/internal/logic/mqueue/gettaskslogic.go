package mqueue

import (
	"context"

	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTasksLogic {
	return &GetTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTasksLogic) GetTasks(req *types.TasksCollectionReq) (resp *types.TasksCollectionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
