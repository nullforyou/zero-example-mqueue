// Code generated by goctl. DO NOT EDIT.
package types

type CreateTaskReq struct {
	CreateTaskItem
}

type CreateTaskResp struct {
	TaskItemResp
}

type TaskItemReq struct {
	TaskName string `path:"task_name"`
}

type CreateTaskItem struct {
	BelongToService string `json:"belong_to_service,optional" comment:"所属服务" validate:"required"`
	CronSpec        string `json:"cron_spec,optional" comment:"周期规格" validate:"required"`
	TaskType        string `json:"task_type,optional" comment:"任务类型" validate:"required"`
	TaskName        string `json:"task_name,optional" comment:"任务名称" validate:"required"`
	TaskRemark      string `json:"task_remark,optional" comment:"任务描述" validate:"required"`
	Target          string `json:"target,optional" comment:"任务目标"`
	State           int    `json:"state,optional" comment:"状态" validate:"omitempty,required,oneof=0 1"`
}

type TaskItemResp struct {
	Id int `json:"id"`
	CreateTaskItem
	UpdatedAt string `json:"updated_at"`
}

type TasksCollectionReq struct {
	BelongToService string `form:"belong_to_service,optional"`
	TaskType        string `form:"task_type,optional"`
	TaskName        string `form:"task_name,optional"`
	State           string `form:"state,optional" validate:"omitempty,required,oneof=enable disable"`
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=10"`
}

type TasksCollectionResp struct {
	Total int            `json:"total"`
	List  []TaskItemResp `json:"list"`
}

type UpdateTaskReq struct {
	CreateTaskItem
}

type UpdateTaskResp struct {
	TaskItemResp
}

type SwitchTaskStateReq struct {
	TaskName string `path:"task_name" validate:"required"`
	State    int    `json:"state,optional" validate:"omitempty,required,oneof=0 1"`
}

type SwitchTaskStateResp struct {
	TaskName string `json:"task_name"`
}

type ExecuteTaskReq struct {
	TaskName string `path:"task_name" validate:"required"`
}

type ExecuteTaskResp struct {
}
