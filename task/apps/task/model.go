package task

import (
	"context"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/tools/ptr"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/webhook"
)

func NewTask(spec TaskSpec) *Task {
	return &Task{
		Id:         uuid.NewString(),
		CreatedAt:  time.Now(),
		TaskSpec:   spec,
		TaskStatus: *NewTaskStatus(),
		Events:     []*event.Event{},
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

	// 执行过程中的事件, 执行日志
	Events []*event.Event `json:"events" gorm:"-" description:"执行过程中的事件"`
}

func (t *Task) TableName() string {
	return "tasks"
}

func (t *Task) String() string {
	return pretty.ToJSON(t)
}

func (t *Task) IsCompleted() bool {
	return slices.Contains([]STATUS{
		STATUS_SUCCESS,
		STATUS_FAILED,
		STATUS_CANCELED,
		STATUS_SKIPPED,
	}, t.Status)
}

func (t *Task) IsRunning() bool {
	return slices.Contains([]STATUS{
		STATUS_QUEUED,
		STATUS_RUNNING,
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

func (t *Task) SetMessage(msg string) *Task {
	t.Message = msg
	return t
}

func (t *Task) WaitComplete(ctx context.Context, scanInterval time.Duration, maxScanCount int) error {
	count := 0
	for {
		if count > maxScanCount {
			return exception.NewInternalServerError("task %s wait complete timeout", t.Id)
		}

		ins, err := GetService().DescribeTask(ctx, NewDescribeTaskRequest(t.Id))
		if err != nil {
			return err
		}
		*t = *ins

		if t.IsCompleted() {
			break
		}

		count++
		log.L().Debug().Msgf("task %s status: %s, retry[%d] next scan after %s", t.Id, t.Status.String(), count, scanInterval)
		time.Sleep(scanInterval)
	}

	return nil
}

func NewTaskSpec(runner string, param *RunParam) *TaskSpec {
	return &TaskSpec{
		Runner:   runner,
		Params:   param,
		Async:    ptr.GetPtr(false),
		Label:    map[string]string{},
		WebHooks: []*webhook.WebHook{},
	}
}

type TaskSpec struct {
	// 任务Id
	TaskId string `json:"task_id" gorm:"-" description:"任务Id, 如果是任务Id,则查询任务执行"`
	// 异步执行时的超时时间
	Timeout string `json:"timeout" gorm:"column:timeout;" description:"异步执行时的超时时间"`
	// 尝试执行,用于做执行前检查
	Async *bool `json:"async" gorm:"column:async;type:bool;" description:"是否是异步任务"`
	// 执行器名称
	Runner string `json:"runner" gorm:"column:type;type:varchar(60);" description:"执行器名称"`
	// 执行器参数
	Params *RunParam `json:"params" gorm:"embedded" description:"任务参数"`
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

func (t *TaskSpec) IsAsync() bool {
	if t.Async == nil {
		return false
	}
	return *t.Async
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

	ctx, fn := context.WithTimeout(context.Background(), timeout)
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

func (t *TaskSpec) SetLabelByMap(m map[string]string) *TaskSpec {
	if t.Label == nil {
		t.Label = map[string]string{}
	}
	maps.Copy(t.Label, m)
	return t
}

func (t *TaskSpec) SetAsync(async bool) *TaskSpec {
	t.Async = &async
	return t
}

func (t *TaskSpec) GetLabel(key string) string {
	if t.Label == nil {
		return ""
	}
	return t.Label[key]
}

type TaskFunc func(ctx context.Context, req any) error

func NewTaskStatus() *TaskStatus {
	return &TaskStatus{
		Status: STATUS_PENDDING,
		Extras: map[string]string{},
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
	// 任务状态是否有变更
	StatusChanged bool `json:"status_changed" gorm:"column:status_changed;type:bool;" description:"任务状态是否有变更"`
	// 失败信息
	Message string `json:"message" gorm:"column:message;type:text;" description:"失败信息"`
	// 失败信息
	Detail string `json:"detail" gorm:"column:detail;type:text;" description:"详情内容"`
	// 任务关联对象
	ReferenceId string `json:"reference_id" gorm:"column:reference_id;type:varchar(200);" description:"关联任务ID"`
	// 管理任务URL
	ReferenceURL string `json:"reference_url" gorm:"column:reference_url;type:varchar(200);" description:"关联任务"`
	// 其他数据
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"任务执行结果" optional:"true"`

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

func (s *TaskStatus) ComputedIsChanged(status STATUS) {
	s.StatusChanged = s.Status == status
}

func (s *TaskStatus) TableName() string {
	return "tasks"
}
