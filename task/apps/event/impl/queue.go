package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/task/apps/event"
)

func (i *EventServiceImpl) RunSaveEventConsumer(ctx context.Context) error {
	err := bus.GetService().QueueSubscribe(ctx, i.EventTopic, func(e *bus.Event) {
		ins := &event.Event{}
		if err := ins.LoadFromEvent(e); err != nil {
			i.log.Error().Msgf("load event error, %s", err)
		}
		if err := datasource.DBFromCtx(ctx).Save(ins).Error; err != nil {
			i.log.Error().Msgf("save event error, %s", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
