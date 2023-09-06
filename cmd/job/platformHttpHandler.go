package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"mqueue/cmd/business"
	"net/http"
	"os"
)

type Scheduler struct {
	BelongToService string `json:"belong_to_service"`
	CronSpec        string `json:"cron_spec"`
	TaskType        string `json:"task_type"`
	TaskName        string `json:"task_name"`
	TaskRemark      string `json:"task_remark"`
	Target          string `json:"target"`
}

type Request struct {
	Node   string `path:"node"`
	ID     int    `form:"id"`
	Header string `header:"X-Header"`
}

func PlatformHttpHandler(ctx context.Context, t *asynq.Task) error {
	var scheduler Scheduler
	if err := json.Unmarshal(t.Payload(), &scheduler); err != nil {
		return fmt.Errorf("<解析Payload失败，payload：%s>: %w", t.Payload(), asynq.SkipRetry)
	}
	logx.Infof("执行【%s】 TaskRemark:%s", business.PlatformHttp, scheduler.TaskRemark)

	return nil
}

func toDo(scheduler Scheduler) {
	req := Request{
		Node:   "foo",
		ID:     1024,
		Header: "foo-header",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, scheduler.Target, req)
	// resp, err := httpc.Do(context.Background(), http.MethodPost, *domain+"/nodes/:node", req)
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(os.Stdout, resp.Body)
}
