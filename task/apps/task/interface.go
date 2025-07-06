package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
)

const (
	APP_NAME = "tasks"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

type Service interface {
	// 任务执行
	Run(context.Context, *TaskSpec) *Task
	// 查询任务列表
	QueryTask(context.Context, *QueryTaskRequest) (*types.Set[*Task], error)
	// 查询任务详情
	DescribeTask(context.Context, *DescribeTaskRequest) (*Task, error)
}

type DescribeTaskRequest struct {
	TaskId string `json:"task_id"`
}

type QueryTaskRequest struct {
}

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
	Id string `json:"id" gorm:"column:id;type:uint;primary_key;" unique:"true" description:"Id"`
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

func (t *Task) Failed() *Task {
	t.Status = STATUS_FAILED
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
		WebHooks: []WebHookSpec{},
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
	WebHooks []WebHookSpec `json:"web_hooks" bson:"web_hooks" gorm:"column:web_hooks;serializer:json;type:json" description:"任务执行结束回调" optional:"true"`

	// 具体的函数
	fn TaskFunc `json:"-" gorm:"-"`
	// 异步任务取消函数
	cancelFn context.CancelFunc `json:"-" gorm:"-"`
}

func (t *TaskSpec) SetAsync(v bool) *TaskSpec {
	t.Async = v
	return t
}

func (t *TaskSpec) BuildTimeoutCtx() context.Context {
	timeout, err := time.ParseDuration(t.Timeout)
	if err != nil {
		timeout = DEFAULT_TIMEOUT
	}

	ctx, fn := context.WithTimeout(context.Background(), timeout)
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
		Events: []event.Event{},
	}
}

type TaskStatus struct {
	// 开始执行时间
	StartAt *time.Time `json:"start_at" gorm:"column:start_at;type:timestamp;" description:"开始执行时间"`
	// 执行结束的时间
	EndAt *time.Time `json:"end_at" gorm:"column:end_at;type:timestamp;" description:"执行结束的时间"`
	// 任务执行状态
	Status STATUS `json:"status" gorm:"column:status;type:tinyint(2);" description:"任务执行状态"`

	// 执行过程中的事件, 执行日志
	Events []event.Event `json:"events" gorm:"column:events;type:json;serializer:json;" description:"执行过程中的事件"`
}

func (s *TaskStatus) SetStartAt(t time.Time) {
	s.StartAt = &t
}

func (s *TaskStatus) SetEndAt(t time.Time) {
	s.EndAt = &t
}

// 任务回调
type WebHookSpec struct {
	// 基本信息
	ID          string `json:"id"`          // WebHook 的唯一标识
	Name        string `json:"name"`        // WebHook 名称
	Description string `json:"description"` // 描述

	// 目标配置
	TargetURL string            `json:"target_url"` // 接收 HTTP 请求的 URL
	Method    string            `json:"method"`     // HTTP 方法，如 POST、GET 等
	Headers   map[string]string `json:"headers"`    // 自定义请求头

	// 触发条件
	Events []string `json:"events"` // 触发的事件类型列表
	Filter string   `json:"filter"` // 事件过滤条件（可选）

	// 安全验证
	Secret          string `json:"secret"`           // 用于签名验证的密钥（HMAC）
	SignatureHeader string `json:"signature_header"` // 签名头的名称，如 "X-Hub-Signature"

	// 请求内容配置
	ContentType string `json:"content_type"` // 请求体类型，如 "application/json"
	Payload     string `json:"payload"`      // 自定义请求体模板（可选）

	// 重试与超时
	RetryCount int `json:"retry_count"` // 失败重试次数
	Timeout    int `json:"timeout"`     // 请求超时时间（毫秒）
}
