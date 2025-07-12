package impl

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/infraboard/modules/task/apps/task"
)

// 处理任务运行队列
func (c *TaskServiceImpl) HandleRunEvents(ctx context.Context) {
	for {
		m, err := c.cancel_reader.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				c.log.Info().Msg("reader closed")
				return
			}
			c.log.Error().Msgf("featch message error, %s", err)
			continue
		}

		// 处理消息
		e := task.NewQueueEvent()
		c.log.Debug().Msgf("message at topic/partition/offset %v/%v/%v", m.Topic, m.Partition, m.Offset)

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err = e.Load(m.Value)
		if err != nil {
			c.log.Error().Msgf("parse event error, %s", err)
			continue
		}

		if e.Type != task.QUEUE_EVENT_TYPE_RUN {
			c.log.Error().Msgf("unknown event type %s", e.Type)
			continue
		}

		// 查询Task实例
		taskIns, err := c.DescribeTask(ctx, task.NewDescribeTaskRequest(e.TaskId))
		if err != nil {
			c.log.Error().Msgf("decribe task error, %s", err)
			continue
		}

		// 处理过的任务不再运行
		if taskIns.Status != task.STATUS_QUEUED {
			c.log.Error().Msgf("任务已经处理: %s", taskIns.Status)
			continue
		}

		// 运行任务
		c.runTask(taskIns)

		// 更新任务状态
		c.updateTaskStatus(ctx, taskIns)
	}
}

func (s *TaskServiceImpl) runTask(ins *task.Task) *task.Task {
	ins.SetStartAt(time.Now())

	switch ins.Type {
	case task.TYPE_FUNCTION:
		fn := ins.GetFn()
		if fn == nil {
			return ins.Failed(fmt.Sprintf("%s fn not found", ins.Id))
		}
		// 执行函数
		ins.Running()
		go func() {
			defer func() {
				ins.Cancel()
				s.RemoveAsyncTask(ins)
			}()
			s.AddAsyncTask(ins)
			if err := fn(ins.BuildTimeoutCtx(), ins.Params); err != nil {
				ins.Failed(err.Error())
			} else {
				ins.Success()
			}
			s.updateTaskStatus(context.Background(), ins)
		}()
	default:
		return ins.Failed(fmt.Sprintf("不支持的类型: %s", ins.Type))
	}

	return ins
}

// 处理任务取消队列
func (c *TaskServiceImpl) HandleCancelEvents(ctx context.Context) {
	for {
		m, err := c.cancel_reader.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF {
				c.log.Info().Msg("reader closed")
				return
			}
			c.log.Error().Msgf("featch message error, %s", err)
			continue
		}

		// 处理消息
		e := task.NewQueueEvent()
		c.log.Debug().Msgf("message at topic/partition/offset %v/%v/%v", m.Topic, m.Partition, m.Offset)

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err = e.Load(m.Value)
		if err != nil {
			c.log.Error().Msgf("parse event error, %s", err)
			continue
		}

		// 查询任务信息
		taskIns, err := c.DescribeTask(ctx, task.NewDescribeTaskRequest(e.TaskId))
		if err != nil {
			c.log.Error().Msgf("decribe cronjob error, %s", err)
			continue
		}

		// 找出任务实例
		taskIns = c.GetAsyncTask(taskIns.Id)
		if taskIns == nil {
			c.log.Info().Msgf("实例[%s]在当前节点[%s]上不存在", e.TaskId, c.node_name)
			continue
		}

		c.log.Info().Msgf("实例[%s]在当前节点[%s]上存在, 正在进行取消处理 ...", e.TaskId, c.node_name)

		// 取消任务, 并移除
		taskIns.Cancel()
		c.RemoveAsyncTask(taskIns)

		// 更新数据裤状态
		taskIns.Status = task.STATUS_CANCELED
		c.updateTaskStatus(ctx, taskIns)
	}
}
