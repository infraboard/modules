package task

import (
	"context"
	"time"
)

const (
	EVENT_LABLE_TASK_ID = "task_id"
)

var (
	DEFAULT_TIMEOUT = time.Second * 30
)

type CONTEXT_TASK_KEY struct{}

func GetTaskFromCtx(ctx context.Context) *Task {
	v := ctx.Value(CONTEXT_TASK_KEY{})
	if v == nil {
		return nil
	}

	return v.(*Task)
}
