package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/webhook"
)

// Run implements webhook.Service.
func (i *WebHookServiceImpl) Run(ctx context.Context, in *webhook.WebHookSpec) *webhook.WebHook {
	hook := webhook.NewWebHook(*in)
	hook.Run(ctx)
	// 保存Hook执行记录
	err := datasource.DBFromCtx(ctx).Save(hook).Error
	if err != nil {
		return hook.Failed(err)
	}
	return hook.Success()
}

// QueryWebHook implements webhook.Service.
func (i *WebHookServiceImpl) QueryWebHook(context.Context, *webhook.QueryWebHookRequest) (*types.Set[*webhook.WebHook], error) {
	return nil, nil
}
