package impl

import (
	"context"

	"github.com/infraboard/modules/maudit/apps/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 存储
// 选择MongoDB
func (i *impl) SaveEvent(ctx context.Context, in *event.EventSet) error {
	i.log.Debug().Msgf("events: %s", in)

	_, err := i.col.InsertMany(ctx, in.ToDocs())
	if err != nil {
		return err
	}
	return nil
}

// 查询
func (i *impl) QueryEvent(ctx context.Context, in *event.QueryEventRequest) (*event.EventSet, error) {
	set := event.NewEventSet()

	filter := bson.M{}

	opt := options.Find()
	opt.SetLimit(int64(in.PageSize))
	opt.SetSkip(in.ComputeOffset())
	cursor, err := i.col.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		e := event.NewEvent()
		if err := cursor.Decode(e); err != nil {
			return nil, err
		}
		set.Add(e)
	}
	return set, nil
}
