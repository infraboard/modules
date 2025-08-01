package runners

import (
	"context"

	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

// 注册Http 调用
func init() {
	task.RegistryRunner(HTTP_CALL, &HttpCall{})
}

const (
	HTTP_CALL = "http_call"
)

type HttpCall struct {
	task.RunnerUnimplemented
}

func (r *HttpCall) Run(ctx context.Context, ins *task.Task) {
	spec := webhook.NewWebHookSpec()

	if err := ins.Params.Load(spec); err != nil {
		ins.Failed(err.Error())
		return
	}

	hook := webhook.NewWebHook(*spec)
	hook.Run(ctx)
	ins.Detail = hook.WebHookStatus.String()

	switch hook.Status {
	case webhook.STATUS_FAILED:
		ins.Status = task.STATUS_FAILED
	case webhook.STATUS_SUCCESS:
		ins.Status = task.STATUS_SUCCESS
	}
}
