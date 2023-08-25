package main

import (
	"gorm.io/plugin/soft_delete"
	"mqueue/cmd/database"
	"time"
)

type Scheduler struct {
	ID              int64                 `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true"`
	BelongToService string                `gorm:"column:belong_to_service;type:varchar(50);not null;default:'';comment:隶属于服务;"`
	CronSpec        string                `gorm:"column:cron_spec;type:varchar(50);not null;default:'';comment:定时任务规格;"`
	TaskType        string                `gorm:"column:task_type;type:varchar(50);not null;default:'';comment:定时任务类型;"`
	TaskName        string                `gorm:"column:task_name;type:varchar(50);not null;default:'';comment:定时任务名称;"`
	TaskRemark      string                `gorm:"column:task_remark;type:varchar(100);not null;default:'';comment:定时任务备注;"`
	Target          string                `gorm:"column:target;type:varchar(200);not null;default:'';comment:目标;"`
	State           int64                 `gorm:"column:state;type:smallint;not null;default:0;comment:状态 1:启用 0:请用;"`
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;type:timestamp"`
	CreatedAt       *time.Time            `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt       *time.Time            `gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (Scheduler) TableName() string {
	return "scheduler"
}

func main() {
	db := database.Connect()
	db.AutoMigrate(&Scheduler{})
}
