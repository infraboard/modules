package task

import (
	"context"
	"encoding/json"
)

var (
	runners = map[string]Runner{}
)

func RegistryRunner(name string, runner Runner) {
	runners[name] = runner
}

func GetRunner(name string) Runner {
	return runners[name]
}

func ListRunner() (ns []string) {
	for k := range runners {
		ns = append(ns, k)
	}
	return
}

// 执行器接口
type Runner interface {
	// 同步模式下: 运行任务，通过ctx取消
	// 异步模式下: 触发任务执行, 任务的状态更新和取消 需要runner单独实现Sync和Cancle
	Run(context.Context, *Task)
	// 异步任务状态同步
	Sync(context.Context, *Task)
	// 异步任务取消
	Cancel(context.Context, *Task)
}

type RunnerUnimplemented struct{}

func (r *RunnerUnimplemented) Run(context.Context, *Task)    {}
func (r *RunnerUnimplemented) Sync(context.Context, *Task)   {}
func (r *RunnerUnimplemented) Cancel(context.Context, *Task) {}

const (
	RUN_PARAM_TYPE_JSON = "json"
)

type RUN_PARAM_TYPE string

func NewJsonRunParam(jsonStr string) *RunParam {
	return &RunParam{
		Type:  RUN_PARAM_TYPE_JSON,
		Value: jsonStr,
	}
}

type RunParam struct {
	// 参数类型
	Type RUN_PARAM_TYPE `json:"type" gorm:"column:run_param_type;type:varchar(60)" description:"参数类型"`
	// 参数值
	Value string `json:"value" gorm:"column:run_param_value;type:text" description:"参数值"`
}

func (r *RunParam) Load(v any) error {
	return json.Unmarshal([]byte(r.Value), v)
}
