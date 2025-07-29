package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/audit"
)

// 读取消息，处理消息
func (c *AuditServiceImpl) RunSaveEventConsumer(ctx context.Context) error {
	bus.GetService().QueueSubscribe(ctx, c.AuditTopic, func(e *bus.Event) {
		ae := audit.NewAuditLog()
		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err := ae.Load(e.Data)
		if err != nil {
			c.log.Error().Msgf("load data error, %s", err)
			return
		}
		if err := datasource.DBFromCtx(ctx).Save(ae).Error; err != nil {
			c.log.Error().Msgf("save event error, %s", err)
		}
	})
	return nil
}
