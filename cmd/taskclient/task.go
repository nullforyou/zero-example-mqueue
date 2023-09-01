// Code generated by goctl. DO NOT EDIT.
// Source: mqueue.proto

package taskclient

import (
	"context"

	"mqueue/cmd/pb/task"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	WorkOrderReply = task.WorkOrderReply
	WorkOrderReq   = task.WorkOrderReq

	Task interface {
		CreateWorkOrder(ctx context.Context, in *WorkOrderReq, opts ...grpc.CallOption) (*WorkOrderReply, error)
	}

	defaultTask struct {
		cli zrpc.Client
	}
)

func NewTask(cli zrpc.Client) Task {
	return &defaultTask{
		cli: cli,
	}
}

func (m *defaultTask) CreateWorkOrder(ctx context.Context, in *WorkOrderReq, opts ...grpc.CallOption) (*WorkOrderReply, error) {
	client := task.NewTaskClient(m.cli.Conn())
	return client.CreateWorkOrder(ctx, in, opts...)
}
