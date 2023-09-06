package mqueue

import (
	"context"
	"errors"
	"go-zero-base/utils/xerr"
	"gorm.io/gorm"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/model"
	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskLogic) UpdateTask(req *types.UpdateTaskReq) (resp *types.UpdateTaskResp, err error) {
	scheduler := model.Scheduler{}
	err = l.svcCtx.DbEngine.Model(scheduler).Where("task_name = ?", req.TaskName).First(&scheduler).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("任务不存在"))
	}
	l.svcCtx.DbEngine.Model(&scheduler).Update("state", req.State).
		Updates(
			model.Scheduler{
				BelongToService: req.BelongToService,
				CronSpec:        req.CronSpec,
				TaskType:        req.TaskType,
				TaskName:        req.TaskName,
				TaskRemark:      req.TaskRemark,
				Target:          req.Target,
			})

	return &types.UpdateTaskResp{
		TaskItemResp: types.TaskItemResp{
			Id: int(scheduler.ID),
			CreateTaskItem: types.CreateTaskItem{
				BelongToService: scheduler.BelongToService,
				CronSpec:        scheduler.CronSpec,
				TaskType:        scheduler.TaskType,
				TaskName:        scheduler.TaskName,
				Target:          scheduler.Target,
				State:           int(scheduler.State),
			},
			UpdatedAt: scheduler.UpdatedAt.Format(business.YYMMDDHHMMSS),
		},
	}, nil
}
