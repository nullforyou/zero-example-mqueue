package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"go-common/tool"
	"mqueue/cmd/business"
	"mqueue/cmd/job"
	"net/http"
	"time"
)

type Scheduler struct {
	BelongToService string `json:"belong_to_service"`
	CronSpec        string `json:"cron_spec"`
	TaskType        string `json:"task_type"`
	TaskName        string `json:"task_name"`
	TaskRemark      string `json:"task_remark"`
	Target          string `json:"target"`
	Payload         string `json:"payload"`
}

type PlatformHttpWorker struct {
	svcCtx *job.ServiceContext
}

func NewPlatformHttpWorker(svcCtx *job.ServiceContext) *PlatformHttpWorker {
	return &PlatformHttpWorker{svcCtx: svcCtx}
}

func (h *PlatformHttpWorker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var scheduler Scheduler
	if err := json.Unmarshal(t.Payload(), &scheduler); err != nil {
		logx.Errorf("[%s]解析Payload失败，payload：%s", business.PlatformHttp, t.Payload())
	}
	logx.Debugf("[%s]执行任务:%s", business.PlatformHttp, scheduler.TaskRemark)
	toDo(h.svcCtx, scheduler)
	return nil
}

func toDo(svcCtx *job.ServiceContext, scheduler Scheduler) {
	request, _ := http.NewRequest("POST", scheduler.Target, bytes.NewBuffer([]byte(scheduler.Payload)))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	request.Header.Set("Authorization", "Bearer "+getJwtToken(svcCtx))

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		logx.Debugf("[%s]执行[%s]任务结果-httpRequest失败：%s", business.PlatformHttp, scheduler.TaskName, err.Error())
		return
	}
	defer response.Body.Close()
	//body, _ := io.ReadAll(response.Body)
	//logx.Debugf("[%s]执行[%s]任务结果-responseStatus:%s;responseBody:%s;", business.PlatformHttp, scheduler.TaskName, response.Status, string(body))
	return
}

func getJwtToken(svcCtx *job.ServiceContext) string {
	key := "mQueue:job-http-jwt"
	exists, _ := svcCtx.Redis.Exists(key)
	if !exists {
		jwtToken, _ := tool.GetJwtToken(svcCtx.Config.Jwt.AccessSecret, svcCtx.Config.Jwt.AccessExpire, 0)
		svcCtx.Redis.SetnxEx(key, jwtToken, 86400)
	}
	jwt, _ := svcCtx.Redis.Get(key)
	return jwt
}
