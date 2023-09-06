package mqueue

import (
	"context"
	"errors"
	"go-zero-base/utils/xerr"
	"gorm.io/gorm"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/query"

	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchTaskStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwitchTaskStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchTaskStateLogic {
	return &SwitchTaskStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwitchTaskStateLogic) SwitchTaskState(req *types.SwitchTaskStateReq) (resp *types.SwitchTaskStateResp, err error) {
	if !(req.State == business.DISABLED || req.State == business.ENABLED) {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorBusiness), xerr.SetMsg("任务状态错误"))
	}

	query.SetDefault(l.svcCtx.DbEngine)
	schedulerDao := query.Scheduler
	schedulerModel, err := schedulerDao.WithContext(context.Background()).Where(schedulerDao.TaskName.Eq(req.TaskName)).First()

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("任务不存在"))
	}
	//如果变更字段为零值，则必须使用这种方法，如果使用结构体会自动过滤
	l.svcCtx.DbEngine.Model(&schedulerModel).Update("state", req.State)

	return &types.SwitchTaskStateResp{TaskName: schedulerModel.TaskName}, nil
}
