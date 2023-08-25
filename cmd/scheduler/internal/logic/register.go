package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

type MqueueScheduler struct {
	ctx       context.Context
	Scheduler *asynq.Scheduler
}

func NewCronScheduler(ctx context.Context, scheduler *asynq.Scheduler) *MqueueScheduler {
	return &MqueueScheduler{
		ctx:       ctx,
		Scheduler: scheduler,
	}
}

type DeferCloseOrderPayload struct {
	OrderSerialNumber string
}

func (l *MqueueScheduler) Register() {
	taskName := "defer:order:close"
	payload, _ := json.Marshal(DeferCloseOrderPayload{OrderSerialNumber: "323222"})
	task := asynq.NewTask(taskName, payload)
	// 每分钟执行一次
	entryID, err := l.Scheduler.Register("@every 1s", task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("定时任务[%s]注册失败:%s", taskName, err)
	}
	logx.WithContext(l.ctx).Debugf("定时任务[%s]注册成功,Id:%s", taskName, entryID)
}
