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
	task.RegistrySyncRunner(DEBUG, &DebugRunner{})
}

const (
	DEBUG = "Debug"
)

type DebugRunner struct {
}

func (r *DebugRunner) Run(ctx context.Context, ins *task.Task) error {
	fmt.Println("执行参数:", ins.Params)
	_, err := event.GetService().AddEvent(ctx, task.NewInfoEvent("开始执行", ins.Id))
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	_, err = event.GetService().AddEvent(ctx, task.NewInfoEvent("执行结束", ins.Id))
	if err != nil {
		return err
	}
	return nil
}
