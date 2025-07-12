package cronjob

import (
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/modules/task/apps/task"
	cron "github.com/robfig/cron/v3"
)

func NewCronJob(spec CronJobSpec) *CronJob {
	return &CronJob{
		Id:          uuid.NewString(),
		CreatedAt:   time.Now(),
		CronJobSpec: spec,
	}
}

type CronJob struct {
	// 任务Id
	Id string `json:"id" gorm:"column:id;type:string;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`

	// Cronjob参数
	CronJobSpec
	// 状态
	CronJobStatus
}

func (c *CronJob) CronEntryID() cron.EntryID {
	return cron.EntryID(c.RefInstanceId)
}

func (c *CronJob) TableName() string {
	return "cronjobs"
}

func NewCronJobSpec(cron string, spec task.TaskSpec) *CronJobSpec {
	return &CronJobSpec{
		Cron:     cron,
		TaskSpec: spec,
	}
}

type CronJobSpec struct {
	// Cron表达式
	Cron string `json:"cron" gorm:"column:cron;type:varchar(120)" description:"Cron表达式"`
	// 是否启用改Cron
	Enabled *bool `json:"enabled" gorm:"column:enabled;type:bool" description:"是否启用改Cron"`
	// Task执行参数
	task.TaskSpec
}

func (t *CronJobSpec) GetEnabled() bool {
	if t.Enabled == nil {
		return false
	}

	return *t.Enabled
}

func (t *CronJobSpec) SetName(name string) *CronJobSpec {
	t.TaskSpec.SetName(name)
	return t
}

func (t *CronJobSpec) SetDescription(desc string) *CronJobSpec {
	t.TaskSpec.SetDescription(desc)
	return t
}

type CronJobStatus struct {
	// 关联的cron实例Id
	RefInstanceId int `json:"ref_instance_id" gorm:"column:ref_instance_id;" description:"关联的cron实例Id"`
	// CronJob执行的Node节点信息
	Node string `json:"node" gorm:"column:node;" description:"CronJob执行的Node节点信息"`
	// 状态更新时间
	UpdateAt *time.Time `json:"update_at" gorm:"column:update_at;type:timestamp;;" description:"状态更新时间"`
	// 状态
	Status STATUS `json:"statsu" gorm:"column:statsu;" description:"CronJob状态"`
}

func (s *CronJobStatus) SetUpdateAt(v time.Time) {
	s.UpdateAt = &v
}

func (s *CronJobStatus) TableName() string {
	return "cronjobs"
}
