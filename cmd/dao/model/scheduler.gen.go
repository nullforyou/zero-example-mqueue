// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameScheduler = "scheduler"

// Scheduler mapped from table <scheduler>
type Scheduler struct {
	ID              int64          `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	BelongToService string         `gorm:"column:belong_to_service;type:varchar(50);not null;comment:隶属于服务" json:"belong_to_service"` // 隶属于服务
	CronSpec        string         `gorm:"column:cron_spec;type:varchar(50);not null;comment:定时任务规格" json:"cron_spec"`                // 定时任务规格
	TaskType        string         `gorm:"column:task_type;type:varchar(50);not null;comment:定时任务类型" json:"task_type"`                // 定时任务类型
	TaskName        string         `gorm:"column:task_name;type:varchar(50);not null;comment:定时任务名称" json:"task_name"`                // 定时任务名称
	TaskRemark      string         `gorm:"column:task_remark;type:varchar(100);not null;comment:定时任务备注" json:"task_remark"`           // 定时任务备注
	Target          string         `gorm:"column:target;type:varchar(200);not null;comment:目标" json:"target"`                         // 目标
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
	CreatedAt       *time.Time     `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

// TableName Scheduler's table name
func (*Scheduler) TableName() string {
	return TableNameScheduler
}
