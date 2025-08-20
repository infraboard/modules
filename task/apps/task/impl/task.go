package impl

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
	"gorm.io/gorm"
)

// 创建任务
func (s *TaskServiceImpl) CreateTask(ctx context.Context, in *task.TaskSpec) (*task.Task, error) {
	ins := task.NewTask(*in)

	err := datasource.DBFromCtx(ctx).Save(ins).Error
	if err != nil {
		return nil, err
	}

	return ins, nil
}

// Run implements task.Service.
func (s *TaskServiceImpl) Run(ctx context.Context, in *task.TaskSpec) (*task.Task, error) {
	var ins *task.Task
	if in.TaskId != "" {
		oldTask, err := s.DescribeTask(ctx, task.NewDescribeTaskRequest(in.TaskId))
		if err != nil {
			return nil, err
		}
		ins = oldTask
	} else {
		ins := task.NewTask(*in)
		// 保存记录
		err := datasource.DBFromCtx(ctx).Save(ins).Error
		if err != nil {
			return nil, err
		}
	}

	// 队列事件
	e := task.NewQueueEvent()
	e.Type = task.QUEUE_EVENT_TYPE_RUN
	e.TaskId = ins.Id

	// 放入运行队列
	err := bus.GetService().Publish(ctx, &bus.Event{
		Subject: s.RunTopic,
		Data:    []byte(e.String()),
	})
	if err != nil {
		return nil, err
	}

	// 更新状态队列中

	return ins, nil
}

// 任务取消
func (s *TaskServiceImpl) Cancel(ctx context.Context, in *task.CancelRequest) (*task.Task, error) {
	ins, err := s.DescribeTask(ctx, task.NewDescribeTaskRequest(in.TaskId))
	if err != nil {
		return nil, err
	}

	e := task.NewQueueEvent()
	e.Type = task.QUEUE_EVENT_TYPE_CANCEL
	e.TaskId = in.TaskId

	// 放入取消队列
	err = bus.GetService().Publish(ctx, &bus.Event{
		Subject: s.CancelTopic,
		Data:    []byte(e.String()),
	})
	if err != nil {
		return nil, err
	}

	ins.Status = task.STATUS_CANCELING
	s.updateTaskStatus(ctx, ins)
	return ins, nil
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
	events, err := event.GetService().QueryEvent(ctx, event.NewQueryEventRequest().
		SetLabel(task.EVENT_LABLE_TASK_ID, ins.Id).
		SetOrderBy(event.ORDER_BY_ASC))
	if err != nil {
		return nil, err
	}
	ins.Events = events.Items
	// 补充WebHook执行数据
	webhooks, err := webhook.GetService().QueryWebHook(ctx, webhook.NewQueryWebHookRequest().SetRefTaskId(ins.Id))
	if err != nil {
		return nil, err
	}
	ins.WebHooks = webhooks.Items

	return ins, nil
}

// QueryTask implements task.Service.
func (i *TaskServiceImpl) QueryTask(ctx context.Context, in *task.QueryTaskRequest) (*types.Set[*task.Task], error) {
	set := types.New[*task.Task]()

	query := datasource.DBFromCtx(ctx).Model(&task.Task{})
	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	err = query.
		Order("created_at desc").
		Offset(int(in.ComputeOffset())).
		Limit(int(in.PageSize)).
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}
	return set, nil
}
