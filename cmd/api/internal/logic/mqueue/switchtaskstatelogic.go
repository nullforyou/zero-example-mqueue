package mqueue

import (
	"context"

	"mqueue/cmd/api/internal/svc"
	"mqueue/cmd/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
