package mqueue

import (
	"context"
	"go-zero-base/utils/xerr"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/model"
	"mqueue/cmd/dao/query"

	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.CreateTaskResp, err error) {
	query.SetDefault(l.svcCtx.DbEngine)
	schedulerDao := query.Scheduler
	count, _ := schedulerDao.WithContext(context.Background()).Where(schedulerDao.TaskName.Eq(req.TaskName)).Count()
	if count > 0 {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorBusiness), xerr.SetMsg("周期任务已存在"))
	}
	schedulerModel := model.Scheduler{
		BelongToService: req.BelongToService,
		CronSpec:        req.CronSpec,
		TaskType:        req.TaskType,
		TaskName:        req.TaskName,
		TaskRemark:      req.TaskRemark,
		Target:          req.Target,
		Payload:         req.Payload,
		State:           int32(req.State),
	}

	result := l.svcCtx.DbEngine.Create(&schedulerModel)

	return &types.CreateTaskResp{
		TaskItemResp: types.TaskItemResp{
			Id: int(schedulerModel.ID),
			CreateTaskItem: types.CreateTaskItem{
				BelongToService: schedulerModel.BelongToService,
				CronSpec:        schedulerModel.CronSpec,
				TaskType:        schedulerModel.TaskType,
				TaskName:        schedulerModel.TaskName,
				Target:          schedulerModel.Target,
				Payload:         schedulerModel.Payload,
				State:           int(schedulerModel.State),
			},
			UpdatedAt: schedulerModel.UpdatedAt.Format(business.YYMMDDHHMMSS),
		},
	}, result.Error
}
