package mqueue

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/model"
	"strings"
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
	var total int64
	queryScheduler := l.svcCtx.DbEngine.Model(model.Scheduler{})
	if !strings.EqualFold(req.BelongToService, "") {
		queryScheduler.Where("belong_to_service = ?", req.BelongToService)
	}
	if !strings.EqualFold(req.TaskType, "") {
		queryScheduler.Where("task_type = ?", req.TaskType)
	}
	if !strings.EqualFold(req.TaskName, "") {
		queryScheduler.Where("task_name like ?", "%"+req.TaskName+"%")
	}
	if strings.EqualFold(req.State, "enable") {
		queryScheduler.Where("state = ?", business.ENABLED)
	}
	if strings.EqualFold(req.State, "disable") {
		queryScheduler.Where("state = ?", business.DISABLED)
	}
	queryScheduler.Count(&total)
	var list []types.TaskItemResp
	err = queryScheduler.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id DESC").Find(&list).Error
	return &types.TasksCollectionResp{
		Total: int(total),
		List:  list,
	}, nil
}
