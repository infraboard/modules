package runners

import (
	"context"
	"fmt"

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

type HttpCall struct{}

func (r *HttpCall) Run(ctx context.Context, req *task.RunParam) (fmt.Stringer, error) {
	spec := webhook.NewWebHookSpec()

	if err := req.Load(spec); err != nil {
		return nil, err
	}

	hook := webhook.NewWebHook(*spec)
	hook.Run(ctx)
	return hook.WebHookStatus, nil
}
