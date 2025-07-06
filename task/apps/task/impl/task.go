package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/task"
	"gorm.io/gorm"
)

// Run implements task.Service.
func (s *TaskServiceImpl) Run(ctx context.Context, in *task.TaskSpec) *task.Task {
	ins := task.NewTask(*in)
	ins.SetStartAt(time.Now())

	// 放数据库
	defer s.saveTask(ctx, ins)

	switch in.Type {
	case task.TYPE_FUNCTION:
		fn := in.GetFn()
		if fn == nil {
			return ins.Failed(fmt.Sprintf("%s fn not found", ins.Id))
		}
		// 执行函数
		if ins.Async {
			ins.Running()
			go func() {
				defer func() {
					in.Cancel()
					s.RemoveAsyncTask(ins)
				}()
				s.AddAsyncTask(ins)
				if err := fn(ins.BuildTimeoutCtx(), ins.Params); err != nil {
					ins.Failed(err.Error())
				} else {
					ins.Success()
				}
				s.updateTask(context.Background(), ins)
			}()
		} else {
			if err := fn(ctx, ins.Params); err != nil {
				return ins.Failed(err.Error())
			}
			ins.Success()
		}
	default:
		return ins.Failed(fmt.Sprintf("不支持的类型: %s", in.Type))
	}
	return ins
}

// DescribeTask implements task.Service.
func (i *TaskServiceImpl) DescribeTask(ctx context.Context, in *task.DescribeTaskRequest) (*task.Task, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &task.Task{}
	if err := query.Where("id = ?", in.TaskId).First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("task %s not found", in.TaskId)
		}
		return nil, err
	}

	// 补充Event数据

	// 补充WebHook执行数据

	return ins, nil
}

// QueryTask implements task.Service.
func (i *TaskServiceImpl) QueryTask(context.Context, *task.QueryTaskRequest) (*types.Set[*task.Task], error) {
	panic("unimplemented")
}
