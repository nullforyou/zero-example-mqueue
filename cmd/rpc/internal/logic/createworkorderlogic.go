package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"go-zero-base/utils/xerr"
	"mqueue/cmd/dao/model"
	"mqueue/cmd/dao/query"
	"time"

	"mqueue/cmd/pb/task"
	"mqueue/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkOrderLogic {
	return &CreateWorkOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWorkOrderLogic) CreateWorkOrder(in *task.WorkOrderReq) (*task.WorkOrderReply, error) {
	return &task.WorkOrderReply{}, nil

	query.SetDefault(l.svcCtx.DbEngine)
	ctx := context.Background()
	scheduler := query.Scheduler
	schedulerModel, err := scheduler.WithContext(ctx).Where(scheduler.TaskName.Eq(in.GetTaskName())).First()
	if err != nil {
		return nil, err
	}
	//插入队列
	err = pushTask(l.svcCtx, schedulerModel)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewBusinessError(xerr.SetCode(xerr.ErrorRpcOther), xerr.SetMsg("插入队列错误")), "把%s插入队列时错误 %+v", schedulerModel.TaskName, err)
	}
	return &task.WorkOrderReply{}, nil
}

/** 临时创建一个队列客户端 */
func pushTask(svcCtx *svc.ServiceContext, schedulerModel *model.Scheduler) error {
	//创建一个客户端
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: svcCtx.Config.Redis.Host, Password: svcCtx.Config.Redis.Pass})

	payload, _ := json.Marshal(schedulerModel)
	newTask := asynq.NewTask(schedulerModel.TaskName, payload)
	//插入队列5秒后处理
	enqueue, err := client.Enqueue(newTask, asynq.ProcessIn(5*time.Second))

	if err != nil {
		return err
	}

	logx.Debugf("enqueued task: id=%s queue=%s", enqueue.ID, enqueue.Queue)

	err = client.Close()
	if err != nil {
		logx.Errorf("把%s插入队列后关闭队列客户端时错误:%+v", schedulerModel.TaskName, err)
	}
	return nil
}
