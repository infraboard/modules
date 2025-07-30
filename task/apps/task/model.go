package task

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/webhook"
)

func NewTask(spec TaskSpec) *Task {
	return &Task{
		Id:         uuid.NewString(),
		CreatedAt:  time.Now(),
		TaskSpec:   spec,
		TaskStatus: *NewTaskStatus(),
	}
}

type Task struct {
	// 任务Id
	Id string `json:"id" gorm:"column:id;type:string;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// 任务定义
	TaskSpec
	// 任务状态
	TaskStatus
}

func (t *Task) TableName() string {
	return "tasks"
}

func (t *Task) IsCompleted() bool {
	return slices.Contains([]STATUS{
		STATUS_SUCCESS,
		STATUS_FAILED,
		STATUS_CANCELED,
		STATUS_SKIPPED,
	}, t.Status)
}

func (t *Task) Failed(msg string) *Task {
	t.Status = STATUS_FAILED
	t.SetEndAt(time.Now())
	t.Message = msg
	return t
}

func (t *Task) Running() *Task {
	t.Status = STATUS_RUNNING
	t.SetEndAt(time.Now())
	return t
}

func (t *Task) Success() *Task {
	t.Status = STATUS_SUCCESS
	t.SetEndAt(time.Now())
	return t
}

func (t *Task) String() string {
	return pretty.ToJSON(t)
}

func NewTaskSpec(runner string, param *RunParam) *TaskSpec {
	return &TaskSpec{
		Runner:   runner,
		Params:   param,
		Label:    map[string]string{},
		WebHooks: []*webhook.WebHook{},
	}
}

type TaskSpec struct {
	// 异步执行时的超时时间
	Timeout string `json:"timeout" gorm:"column:timeout;" description:"异步执行时的超时时间"`
	// 执行器名称
	Runner string `json:"runner" gorm:"column:type;type:varchar(60);" description:"执行器名称"`
	// 执行器参数
	Params *RunParam `json:"params" gorm:"column:params;serializer:json;type:json" description:"任务参数"`
	// 任务名称
	Name string `json:"name" gorm:"column:name;type:varchar(200);" description:"任务名称"`
	// 任务名称
	Description string `json:"description" gorm:"column:description;type:text;" description:"任务描述"`
	// 尝试执行,用于做执行前检查
	DryRun bool `json:"dryrun" gorm:"column:dryrun;type:bool;" description:"尝试执行,用于做执行前检查"`
	// 事件所属资源
	Resource string `json:"resource" gorm:"column:resource;type:varchar(120);" description:"事件所属资源"`
	// 任务标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json;type:json" description:"任务标签" optional:"true"`

	// 任务执行结束回调
	WebHooks []*webhook.WebHook `json:"web_hooks" bson:"web_hooks" gorm:"column:web_hooks;serializer:json;type:json" description:"任务执行结束回调" optional:"true"`
}

func (t *TaskSpec) SetName(name string) *TaskSpec {
	t.Name = name
	return t
}

func (t *TaskSpec) SetDescription(desc string) *TaskSpec {
	t.Description = desc
	return t
}

func (t *TaskSpec) AddWebHook(hs ...*webhook.WebHook) *TaskSpec {
	t.WebHooks = append(t.WebHooks, hs...)
	return t
}

// 注入上下文当中
func (t *Task) BuildTimeoutCtx() context.Context {
	timeout, err := time.ParseDuration(t.Timeout)
	if err != nil {
		timeout = DEFAULT_TIMEOUT
	}

	ctx := context.WithValue(context.Background(), CONTEXT_TASK_KEY{}, t)
	ctx, fn := context.WithTimeout(ctx, timeout)
	t.cancelFn = fn
	return ctx
}

func (t *Task) Cancel() {
	if t.cancelFn != nil {
		t.cancelFn()
	}
}

func (t *TaskSpec) SetLabel(key, value string) *TaskSpec {
	if t.Label == nil {
		t.Label = map[string]string{}
	}
	t.Label[key] = value
	return t
}

type TaskFunc func(ctx context.Context, req any) error

func NewTaskStatus() *TaskStatus {
	return &TaskStatus{
		Status: STATUS_QUEUED,
		Events: []*event.Event{},
	}
}

type TaskStatus struct {
	// 开始执行时间
	StartAt *time.Time `json:"start_at" gorm:"column:start_at;type:timestamp;" description:"开始执行时间"`
	// 执行结束的时间
	EndAt *time.Time `json:"end_at" gorm:"column:end_at;type:timestamp;" description:"执行结束的时间"`
	// 任务状态更新时间
	UpdateAt *time.Time `json:"update_at" gorm:"column:update_at;type:timestamp;" description:"任务状态更新时间"`
	// 任务执行状态
	Status STATUS `json:"status" gorm:"column:status;type:tinyint(2);" description:"任务执行状态"`
	// 失败信息
	Message string `json:"message" gorm:"column:message;type:text;" description:"失败信息"`
	// 失败信息
	Detail string `json:"detail" gorm:"column:detail;type:text;" description:"详情内容"`

	// 执行过程中的事件, 执行日志
	Events []*event.Event `json:"events" gorm:"column:events;type:json;serializer:json;" description:"执行过程中的事件"`

	// 异步任务取消函数
	cancelFn context.CancelFunc `json:"-" gorm:"-"`
}

func (s *TaskStatus) SetStartAt(t time.Time) {
	s.StartAt = &t
}

func (s *TaskStatus) SetEndAt(t time.Time) {
	s.EndAt = &t
}

func (s *TaskStatus) SetUpdateAt(t time.Time) {
	s.UpdateAt = &t
}

func (s *TaskStatus) TableName() string {
	return "tasks"
}
