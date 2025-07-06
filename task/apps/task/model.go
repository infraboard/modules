package task

import (
	"context"
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

func NewFnTask(fn TaskFunc, params any) *TaskSpec {
	return &TaskSpec{
		Type:     TYPE_FUNCTION,
		fn:       fn,
		Params:   params,
		Label:    map[string]string{},
		WebHooks: []*webhook.WebHook{},
	}
}

type TaskSpec struct {
	// 是否异步执行
	Async bool `json:"async" gorm:"column:async;" description:"是否异步执行"`
	// 异步执行时的超时时间
	Timeout string `json:"timeout" gorm:"column:timeout;" description:"异步执行时的超时时间"`
	// 任务类型
	Type TYPE `json:"type" gorm:"column:id;type:varchar(60);" description:"任务类型"`
	// 任务的参数
	Params any `json:"params" gorm:"column:params;serializer:json;type:json" description:"任务参数"`
	// 任务标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json;type:json" description:"任务标签" optional:"true"`

	// 任务执行结束回调
	WebHooks []*webhook.WebHook `json:"web_hooks" bson:"web_hooks" gorm:"column:web_hooks;serializer:json;type:json" description:"任务执行结束回调" optional:"true"`

	// 具体的函数
	fn TaskFunc `json:"-" gorm:"-"`
	// 异步任务取消函数
	cancelFn context.CancelFunc `json:"-" gorm:"-"`
}

func (t *TaskSpec) SetAsync(v bool) *TaskSpec {
	t.Async = v
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

func (t *TaskSpec) Cancel() {
	if t.cancelFn != nil {
		t.cancelFn()
	}
}

func (t *TaskSpec) GetFn() TaskFunc {
	return t.fn
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

	// 执行过程中的事件, 执行日志
	Events []*event.Event `json:"events" gorm:"column:events;type:json;serializer:json;" description:"执行过程中的事件"`
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
