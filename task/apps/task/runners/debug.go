package runners

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/task"
)

// 注册DebugRunner 用于调试
func init() {
	task.RegistryRunner(DEBUG, &DebugRunner{})
}

const (
	DEBUG = "Debug"
)

type DebugRunner struct{}

func (r *DebugRunner) Run(ctx context.Context, req *task.RunParam) (fmt.Stringer, error) {
	fmt.Println(req.Value)

	ins := task.GetTaskFromCtx(ctx)

	_, err := event.GetService().AddEvent(ctx, task.NewInfoEvent("开始执行", ins.Id))
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	_, err = event.GetService().AddEvent(ctx, task.NewInfoEvent("执行结束", ins.Id))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
