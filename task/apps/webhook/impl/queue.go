package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/task/apps/webhook"
)

func (i *WebHookServiceImpl) RunWebHookConsumer(ctx context.Context) error {
	err := bus.GetService().QueueSubscribe(ctx, i.WebhookQueue, func(e *bus.Event) {
		hook, err := i.DescribeWebHook(ctx, webhook.NewDescribeWebHookRequest(string(e.Data)))
		if err != nil {
			i.log.Error().Msgf("describe hook error, %s", err)
			return
		}

		// 运行
		hook.Run(ctx)

		// 更新状态
		if err := datasource.DBFromCtx(ctx).Where("id = ?", hook.Id).Updates(hook.WebHookStatus).Error; err != nil {
			i.log.Error().Msgf("save event error, %s", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
