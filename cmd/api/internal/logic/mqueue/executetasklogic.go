package mqueue

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-base/utils/xerr"
	"gorm.io/gorm"
	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"
	"mqueue/cmd/dao/model"
	"mqueue/cmd/taskclient"
)

type ExecuteTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteTaskLogic {
	return &ExecuteTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteTaskLogic) ExecuteTask(req *types.ExecuteTaskReq) (resp *types.ExecuteTaskResp, err error) {
	scheduler := model.Scheduler{}
	err = l.svcCtx.DbEngine.Model(scheduler).Where("task_name = ?", req.TaskName).First(&scheduler).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("任务不存在"))
	}
	_, err = l.svcCtx.MQueueRpc.CreateWorkOrder(l.ctx, &taskclient.WorkOrderReq{TaskName: req.TaskName})
	if err != nil {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorBusiness), xerr.SetMsg("入队列时错误"))
	}
	return &types.ExecuteTaskResp{}, nil
}
