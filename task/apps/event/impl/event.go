package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
	"gorm.io/datatypes"
)

// AddEvent implements event.Service.
func (i *EventServiceImpl) AddEvent(ctx context.Context, in *event.EventSpec) (*event.Event, error) {
	ins := event.NewEvent(*in)
	err := bus.GetService().Publish(ctx, ins.ToBusEvent(i.EventTopic))
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *EventServiceImpl) SaveEvent(ctx context.Context, ins *event.Event) error {
	err := bus.GetService().Queue(ctx, i.EventTopic, func(e *bus.Event) {
		if err := datasource.DBFromCtx(ctx).Save(ins).Error; err != nil {
			i.log.Error().Msgf("save event error, %s", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// QueryEvent implements event.Service.
func (i *EventServiceImpl) QueryEvent(ctx context.Context, in *event.QueryEventRequest) (*types.Set[*event.Event], error) {
	set := types.NewSet[*event.Event]()

	query := datasource.DBFromCtx(ctx).Model(&event.Event{})
	if in.OrderBy != "" {
		query = query.Order(fmt.Sprintf("time %s", in.OrderBy))
	}

	for key, value := range in.Label {
		query = query.Where(datatypes.JSONQuery("label").Equals(value, key))
	}

	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	err = query.
		Offset(int(in.ComputeOffset())).
		Limit(int(in.PageSize)).
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}

	return set, nil
}
