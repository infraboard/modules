package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/task"
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
			// 需要保存事件
			e := task.NewErrorEvent("fn not found", ins.Id)
			s.saveEvent(ctx, e)
			return ins.Failed(e.First().Message)
		}
		// 执行函数
		if in.Async {
			go func() {
				defer func() {
					in.Cancel()
					s.RemoveAsyncTask(ins)
				}()
				s.AddAsyncTask(ins)
				if err := fn(in.BuildTimeoutCtx(), ins.Params); err != nil {
					s.saveEvent(ctx, task.NewErrorEvent(err.Error(), ins.Id))
				}
			}()
		} else {
			if err := fn(ctx, in.Params); err != nil {
				s.saveEvent(ctx, task.NewErrorEvent(err.Error(), ins.Id))
				return ins.Failed(err.Error())
			}
		}

		ins.Success()
	default:
		return ins.Failed(fmt.Sprintf("不支持的类型: %s", in.Type))
	}
	return ins
}

// DescribeTask implements task.Service.
func (i *TaskServiceImpl) DescribeTask(context.Context, *task.DescribeTaskRequest) (*task.Task, error) {
	panic("unimplemented")
}

// QueryTask implements task.Service.
func (i *TaskServiceImpl) QueryTask(context.Context, *task.QueryTaskRequest) (*types.Set[*task.Task], error) {
	panic("unimplemented")
}
