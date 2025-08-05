package impl

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
	"github.com/infraboard/mcube/v2/ioc/config/mcron"
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/infraboard/modules/task/apps/task"
)

// 处理CronJob更新队列
func (c *CronJobServiceImpl) HandleUpdateEvents(ctx context.Context) {
	c.log.Info().Msgf("start handle cronjob update events ...")
	err := bus.GetService().TopicSubscribe(ctx, c.UpdateTopic, func(e *bus.Event) {
		// 处理消息
		te := cronjob.NewQueueEvent()

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		err := te.Load(e.Data)
		if err != nil {
			c.log.Error().Msgf("parse event error, %s", err)
			return
		}

		// 查询
		job, err := c.DescribeCronJob(ctx, cronjob.NewDescribeCronJobRequest(te.CronJobId))
		if err != nil {
			c.log.Error().Msgf("decribe cronjob error, %s", err)
			return
		}

		// 处理事件
		switch te.Type {
		case cronjob.QUEUE_EVENT_TYPE_ADD:
			c.addCronjob(job)
		case cronjob.QUEUE_EVENT_TYPE_DELETE:
			c.deleteCronJob(ctx, job)
		case cronjob.QUEUE_EVENT_TYPE_UPDATE:
			c.updateCronjob(ctx, job)
		default:
			c.log.Error().Msgf("unknown event type %s", te.Type)
		}
	})
	if err != nil {
		c.log.Error().Msgf("subscribe update event error, %s", err)
	}
}

func (c *CronJobServiceImpl) deleteCronJob(ctx context.Context, job *cronjob.CronJob) {
	// 找到该实例
	entry := mcron.Get().Entry(job.CronEntryID())
	if entry.ID == 0 {
		c.log.Info().Msgf("实例在当前节点[%s]上不存在", c.node_name)
		return
	}

	mcron.Get().Remove(entry.ID)
	c.log.Info().Msgf("删除当前节点[%s]上的实例: %d", c.node_name, entry.ID)
	err := datasource.DBFromCtx(ctx).Where("id = ?", job.Id).Delete(job).Error
	if err != nil {
		c.log.Error().Msgf("删除数据库记录异常: %s", err)
		return
	}
}

func (c *CronJobServiceImpl) updateCronjob(ctx context.Context, job *cronjob.CronJob) {
	// 找到该实例
	entry := mcron.Get().Entry(job.CronEntryID())
	if entry.ID == 0 {
		c.log.Info().Msgf("实例在当前节点[%s]上不存在", c.node_name)
		return
	}

	// 有更新 则删除后, 重新添加
	mcron.Get().Remove(entry.ID)
	c.log.Info().Msgf("删除当前节点[%s]上的实例: %d", c.node_name, entry.ID)

	// 才真正的创建cron 关联cron实例
	if job.GetEnabled() {
		refId, err := mcron.RunAndAddFunc(job.Cron, func() {
			task.GetService().Run(context.Background(), &job.TaskSpec)
		})
		if err != nil {
			c.log.Error().Msgf("重新添加更新后的cron失败: %s", err)
			return
		}
		job.RefInstanceId = int(refId)
		job.Node = c.node_name
	} else {
		c.log.Info().Msgf("cron未启用")
	}

	job.SetUpdateAt(time.Now())
	if err := datasource.DBFromCtx(ctx).Where("id = ?", job.Id).Updates(job.CronJobStatus).Error; err != nil {
		c.log.Error().Msgf("重新添加更新后的cron失败: %s", err)
		return
	}
}

func (c *CronJobServiceImpl) addCronjob(job *cronjob.CronJob) {
	refId, err := mcron.Get().AddFunc(job.Cron, func() {
		// 需要加锁执行，并且判断时间间隔
		m := lock.L().New(job.Id, c.taskLockTTL)
		if err := m.TryLock(c.ctx); err != nil {
			c.log.Error().Msg(err.Error())
			return
		}
		defer m.UnLock(c.ctx)

		// 查询Job最新一次执行时间
		jobIns, err := c.DescribeCronJob(c.ctx, cronjob.NewDescribeCronJobRequest(job.Id))
		if err != nil {
			c.log.Error().Msg(err.Error())
			return
		}

		// 执行间隔判断
		// jobIns.LatestRunAt

		t, err := task.GetService().Run(context.Background(), &jobIns.TaskSpec)
		if err != nil {
			c.log.Error().Msg(err.Error())
			return
		}

		c.log.Debug().Msg(t.String())
		// 执行完成 需要更新cron最新一次执行时间
	})
	if err != nil {
		job.Failed(err.Error())
		return
	}
	job.RefInstanceId = int(refId)
	job.Node = c.node_name
}
