package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/audit"
	"gorm.io/datatypes"
)

// 存储
func (i *AuditServiceImpl) SaveEvent(ctx context.Context, in *audit.AuditLog) error {
	return bus.GetService().Publish(ctx, in.ToBusEvent(i.AuditTopic))
}

// 查询
func (i *AuditServiceImpl) QueryEvent(ctx context.Context, in *audit.QueryAuditLogRequest) (*types.Set[*audit.AuditLog], error) {
	set := types.NewSet[*audit.AuditLog]()

	query := datasource.DBFromCtx(ctx).Model(&audit.AuditLog{})

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
