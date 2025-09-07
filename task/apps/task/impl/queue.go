package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/modules/task/apps/task"
)

// 处理任务运行队列
func (c *TaskServiceImpl) HandleRunEvents(ctx context.Context) {
	c.log.Info().Msgf("start handle task run events ...")

	err := bus.GetService().TopicSubscribe(ctx, c.RunTopic, func(e *bus.Event) {
		// 处理消息
		te := task.NewQueueEvent()

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err := te.Load(e.Data)
		if err != nil {
			c.log.Error().Msgf("parse event error, %s", err)
			return
		}

		if te.Type != task.QUEUE_EVENT_TYPE_RUN {
			c.log.Error().Msgf("unknown event type %s", te.Type)
			return
		}

		// 查询Task实例
		taskIns, err := c.DescribeTask(ctx, task.NewDescribeTaskRequest(te.TaskId))
		if err != nil {
			c.log.Error().Msgf("decribe task error, %s", err)
			return
		}

		// 判断任务是否在排队的过程中取消了
		if taskIns.Status == task.STATUS_CANCELED {
			c.log.Error().Msgf("任务已经取消: %s", taskIns.Status)
			return
		}

		c.log.Info().Msgf("[开始]开始在节点[%s]上异步执行task: %s ...", c.node_name, taskIns.Id)
		// 运行任务
		c.runTask(taskIns)
		// 更新任务状态
		c.updateTaskStatus(ctx, taskIns)

		c.log.Info().Msgf("[结束]在节点[%s]上异步执行task: %s", c.node_name, taskIns.Id)
	})
	if err != nil {
		c.log.Error().Msgf("subscribe run event error, %s", err)
	}
}

func (s *TaskServiceImpl) runTask(ins *task.Task) *task.Task {
	ins.SetStartAt(time.Now())
	ins.Status = task.STATUS_RUNNING

	if ins.Async {
		// 异步执行
		runner := task.GetAsyncRunner(ins.Runner)
		if runner == nil {
			return ins.Failed(fmt.Sprintf("runner %s not found", ins.Runner))
		}

		// 只是触发成功, 任务状态还是运行中
		err := runner.Start(ins.BuildTimeoutCtx(), ins)
		if err != nil {
			return ins.Failed(err.Error())
		}
	} else {
		// 同步执行, 是一个对象，在内存中持续运行
		runner := task.GetSyncRunner(ins.Runner)
		if runner == nil {
			return ins.Failed(fmt.Sprintf("runner %s not found", ins.Runner))
		}

		// 执行函数
		s.AddAsyncTask(ins)
		defer s.RemoveAsyncTask(ins)

		err := runner.Run(ins.BuildTimeoutCtx(), ins)
		if err != nil {
			return ins.Failed(err.Error())
		}
		ins.Success().SetMessage("任务运行成功")
	}
	return ins
}

// 处理任务取消队列
func (c *TaskServiceImpl) HandleCancelEvents(ctx context.Context) {
	c.log.Info().Msgf("start handle task cancel events ...")
	err := bus.GetService().TopicSubscribe(ctx, c.CancelTopic, func(e *bus.Event) {
		// 处理消息
		te := task.NewQueueEvent()

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err := te.Load(e.Data)
		if err != nil {
			c.log.Error().Msgf("parse event error, %s", err)
			return
		}

		// 查询任务信息
		taskIns, err := c.DescribeTask(ctx, task.NewDescribeTaskRequest(te.TaskId))
		if err != nil {
			c.log.Error().Msgf("decribe cronjob error, %s", err)
			return
		}

		// 找出任务实例
		taskIns = c.GetAsyncTask(taskIns.Id)
		if taskIns == nil {
			c.log.Info().Msgf("实例[%s]在当前节点[%s]上不存在", te.TaskId, c.node_name)
			return
		}

		c.log.Info().Msgf("实例[%s]在当前节点[%s]上存在, 正在进行取消处理 ...", te.TaskId, c.node_name)

		// 取消任务, 并移除
		taskIns.Cancel()
		c.RemoveAsyncTask(taskIns)

		// 更新数据裤状态
		taskIns.Status = task.STATUS_CANCELED
		c.updateTaskStatus(ctx, taskIns)
	})
	if err != nil {
		c.log.Error().Msgf("subscribe cancel event error, %s", err)
	}
}
