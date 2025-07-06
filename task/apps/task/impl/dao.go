package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/task/apps/task"
)

func (s *TaskServiceImpl) save(ctx context.Context, ins *task.Task) {
	err := datasource.DBFromCtx(ctx).Save(ins).Error
	if err != nil {
		s.log.Error().Msgf("save task error, %s", err)
	}
}
