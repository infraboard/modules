package impl_test

import (
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

func TestAddCronJob(t *testing.T) {
	req := task.NewTaskSpec(task.DEBUG_RUNNER, task.NewJsonRunParam("test"))

	req.AddWebHook(webhook.NewWebHook(webhook.WebHookSpec{
		TargetURL:  "https://www.baidu.com/",
		Method:     "GET",
		Conditions: task.StatusCompleteString(),
	}))

	ins, err := svc.AddCronJob(t.Context(), cronjob.NewCronJobSpec("@every 3s", *req).SetName("cronjob测试"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)

	time.Sleep(10 * time.Second)
}
