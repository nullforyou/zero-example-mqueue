package mqueue

import (
	"context"
	"errors"
	"go-zero-base/utils/xerr"
	"gorm.io/gorm"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/model"

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
	scheduler := model.Scheduler{}
	err = l.svcCtx.DbEngine.Model(scheduler).Where("task_name = ?", req.TaskName).First(&scheduler).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("任务不存在"))
	}

	return &types.TaskItemResp{
		CreateTaskItem: types.CreateTaskItem{
			BelongToService: scheduler.BelongToService,
			CronSpec:        scheduler.CronSpec,
			TaskType:        scheduler.TaskType,
			TaskName:        scheduler.TaskName,
			TaskRemark:      scheduler.TaskRemark,
			Target:          scheduler.Target,
			State:           int(scheduler.State),
		},
		UpdatedAt: scheduler.UpdatedAt.Format(business.YYMMDDHHMMSS),
	}, err
}
