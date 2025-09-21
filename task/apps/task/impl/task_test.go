package impl_test

import (
	"testing"
	"time"

	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/task/runners"
	"github.com/infraboard/modules/task/apps/webhook"
)

func TestRun(t *testing.T) {
	req := task.NewTaskSpec(runners.DEBUG, task.NewJsonRunParam("test")).SetLabel("deploy_id", "dpl-01")

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

	time.Sleep(time.Second * 10)
}

const (
	TASK_ID = "f15976a8-42ec-42d6-89cb-73b149f72ac9"
)

func TestDescribeTask(t *testing.T) {
	ins, err := svc.DescribeTask(t.Context(), task.NewDescribeTaskRequest(TASK_ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestQueryTask(t *testing.T) {
	req := task.NewQueryTaskRequest().SetLabel("deploy_id", "dpl-01")
	ins, err := svc.QueryTask(t.Context(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestUpdateTaskStatus(t *testing.T) {
	ins, err := svc.UpdateTaskStatus(t.Context(), task.NewUpdateTaskStatusRequest(TASK_ID))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
