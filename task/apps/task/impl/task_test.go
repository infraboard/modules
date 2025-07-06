package impl_test

import (
	"context"
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

func TestRun(t *testing.T) {
	fnTask := task.NewFnTask(func(ctx context.Context, req any) error {
		t.Log(req.(map[string]string)["param01"])

		ins := task.GetTaskFromCtx(ctx)

		_, err := event.GetService().AddEvent(t.Context(), task.NewInfoEvent("开始执行", ins.Id))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(3 * time.Second)
		_, err = event.GetService().AddEvent(t.Context(), task.NewInfoEvent("执行结束", ins.Id))
		if err != nil {
			t.Fatal(err)
		}

		// secrt同步
		return nil
	}, map[string]string{
		"param01": "01",
	})

	fnTask.SetAsync(true)
	fnTask.AddWebHook(*webhook.NewWebHook(webhook.WebHookSpec{
		TargetURL:  "https://www.baidu.com/",
		Method:     "GET",
		Conditions: task.StatusCompleteString(),
	}))

	ins := svc.Run(t.Context(), fnTask)
	t.Log(ins)

	time.Sleep(5 * time.Second)
	ins, err := svc.DescribeTask(t.Context(), task.NewDescribeTaskRequest(ins.Id))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

const (
	TASK_ID = "f74f2d7e-1926-4f4e-bff3-f4766af34ddb"
)

func TestDescribeTask(t *testing.T) {
	ins, err := svc.DescribeTask(t.Context(), task.NewDescribeTaskRequest(TASK_ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
