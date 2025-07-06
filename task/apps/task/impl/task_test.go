package impl_test

import (
	"context"
	"testing"

	"github.com/infraboard/modules/task/apps/task"
)

func TestRun(t *testing.T) {
	fnTask := task.NewFnTask(func(ctx context.Context, req any) error {
		t.Log(req.(map[string]string)["param01"])
		// secrt同步
		return nil
	}, map[string]string{
		"param01": "01",
	})

	fnTask.SetAsync(true)
	resp := svc.Run(t.Context(), fnTask)
	t.Log(resp)
}
