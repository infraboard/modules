package impl_test

import (
	"context"
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/task"
)

func TestAddCronJob(t *testing.T) {
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

	ins, err := svc.AddCronJob(t.Context(), cronjob.NewCronJobSpec("@every 3s", *fnTask).SetName("cronjob测试"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	time.Sleep(10 * time.Second)
}
