package cronjob

import (
	"time"

	"github.com/infraboard/modules/task/apps/task"
)

type CronJob struct {
	// 任务Id
	Id string `json:"id" gorm:"column:id;type:string;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// Cronjob参数
	CronJobSpec
}

func (c *CronJob) TableName() string {
	return "cronjobs"
}

type CronJobSpec struct {
	// Cron表达式
	Cron string `json:"cron" gorm:"column:cron;type:varchar(120)" description:"Cron表达式"`
	// Task执行参数
	task.TaskSpec
}
