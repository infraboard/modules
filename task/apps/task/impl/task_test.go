package impl_test

import (
	"testing"

	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

func TestRun(t *testing.T) {
	req := task.NewTaskSpec(task.DEBUG_RUNNER, task.NewJsonRunParam("test"))

	req.AddWebHook(webhook.NewWebHook(webhook.WebHookSpec{
		TargetURL:  "https://www.baidu.com/",
		Method:     "GET",
		Conditions: task.StatusCompleteString(),
	}))

	ins, err := svc.Run(t.Context(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

const (
	TASK_ID = "11029dbd-22cd-4c14-9c09-249ff2d84212"
)

func TestDescribeTask(t *testing.T) {
	ins, err := svc.DescribeTask(t.Context(), task.NewDescribeTaskRequest(TASK_ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestQueryTask(t *testing.T) {
	ins, err := svc.QueryTask(t.Context(), task.NewQueryTaskRequest())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
