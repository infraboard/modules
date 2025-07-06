package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
	"gorm.io/gorm"
)

// AddEvent implements event.Service.
func (i *EventServiceImpl) AddEvent(ctx context.Context, in *types.Set[*event.EventSpec]) (*types.Set[*event.Event], error) {
	errors := []string{}

	// 构建对象
	events := types.NewSet[*event.Event]()
	in.ForEach(func(spec *event.EventSpec) {
		events.Add(event.NewEvent(*spec))
	})

	err := datasource.DBFromCtx(ctx).Transaction(func(tx *gorm.DB) error {
		events.ForEach(func(e *event.Event) {
			err := tx.Save(e).Error
			if err != nil {
				errors = append(errors, err.Error())
			}
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return events, fmt.Errorf("%s", strings.Join(errors, ","))
	}

	return events, nil
}

// QueryEvent implements event.Service.
func (i *EventServiceImpl) QueryEvent(ctx context.Context, in *event.QueryEventRequest) (*types.Set[*event.Event], error) {
	panic("unimplemented")
}
