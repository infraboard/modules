package task

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/modules/task/apps/event"
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
	Run(context.Context, *RunParam) (fmt.Stringer, error)
}

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
	Type  RUN_PARAM_TYPE `json:"type"`
	Value any            `json:"value"`
}

// 注册DebugRunner 用于调试
func init() {
	RegistryRunner(DEBUG_RUNNER, &DebugRunner{})
}

const (
	DEBUG_RUNNER = "Debug"
)

type DebugRunner struct{}

func (r *DebugRunner) Run(ctx context.Context, req *RunParam) (fmt.Stringer, error) {
	fmt.Println(req.Value)

	ins := GetTaskFromCtx(ctx)

	_, err := event.GetService().AddEvent(ctx, NewInfoEvent("开始执行", ins.Id))
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	_, err = event.GetService().AddEvent(ctx, NewInfoEvent("执行结束", ins.Id))
	if err != nil {
		return nil, err
	}

	// secrt同步
	return nil, nil
}
