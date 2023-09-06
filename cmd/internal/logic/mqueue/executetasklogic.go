package mqueue

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-base/utils/xerr"
	"gorm.io/gorm"
	"mqueue/cmd/business"
	"mqueue/cmd/dao/model"
	"mqueue/cmd/internal/svc"
	"mqueue/cmd/internal/types"
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

	err = pushTask(l.svcCtx, &scheduler)

	if err != nil {
		return nil, errors.Wrapf(xerr.NewBusinessError(xerr.SetCode(xerr.ErrorRpcOther), xerr.SetMsg("插入队列错误")), "把%s插入队列时错误 %+v", scheduler.TaskName, err)
	}
	return &types.ExecuteTaskResp{}, nil
}

/** 临时创建一个队列客户端 */
func pushTask(svcCtx *svc.ServiceContext, schedulerModel *model.Scheduler) error {
	//创建一个客户端
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: svcCtx.Config.Redis.Host, Password: svcCtx.Config.Redis.Pass})

	payload, _ := json.Marshal(schedulerModel)
	newTask := asynq.NewTask(business.PlatformHttp, payload, asynq.Queue("critical"), asynq.MaxRetry(1))
	//插入队列5秒后处理
	enqueue, err := client.Enqueue(newTask)

	if err != nil {
		return err
	}

	logx.Debugf("已进入队列: id=%s queue=%s", enqueue.ID, enqueue.Queue)

	err = client.Close()
	if err != nil {
		logx.Errorf("把%s插入队列后关闭队列客户端时错误:%+v", schedulerModel.TaskName, err)
	}
	return nil
}
