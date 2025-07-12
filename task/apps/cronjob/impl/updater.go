package impl

import (
	"context"
	"io"
	"time"
)

// 读取消息，处理消息
func (c *CronJobServiceImpl) Run(ctx context.Context) error {
	for {
		m, err := c.updater.FetchMessage(ctx)
		if err != nil {
			if err == io.EOF {
				c.log.Info().Msg("reader closed")
				return nil
			}
			c.log.Error().Msgf("featch message error, %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 处理消息
		// e := event.NewEvent()
		c.log.Debug().Msgf("message at topic/partition/offset %v/%v/%v", m.Topic, m.Partition, m.Offset)

		// 发送的数据时Json格式, 接收用的JSON, 发送也需要使用JSON
		// err = e.Load(m.Value)
		// if err == nil {
		// 	if err := event.GetService().SaveEvent(ctx, event.NewEventSet().Add(e)); err != nil {
		// 		c.log.Error().Msgf("save event error, %s", err)
		// 	}
		// }

		// 处理完消息后需要提交该消息已经消费完成, 消费者挂掉后保存消息消费的状态
		if err := c.updater.CommitMessages(ctx, m); err != nil {
			c.log.Error().Msgf("failed to commit messages: %s", err)
		}
	}
}
