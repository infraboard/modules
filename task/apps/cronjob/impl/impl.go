package impl

import (
	"context"
	"os"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	ioc_kafka "github.com/infraboard/mcube/v2/ioc/config/kafka"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

func init() {
	ioc.Controller().Registry(&CronJobServiceImpl{
		EnableUpdate: true,
		UpdateTopic:  "cronjob_update_events",
	})
}

var _ cronjob.Service = (*CronJobServiceImpl)(nil)

type CronJobServiceImpl struct {
	ioc.ObjectImpl

	log      *zerolog.Logger
	hostname string

	// 更新队列(cron更新, 删除, 启用, 禁用)
	updater *kafka.Reader
	// 允许时上下文
	ctx context.Context

	// 启用更新队列
	EnableUpdate bool `toml:"enable_update" json:"enable_update" yaml:"enable_update"  env:"ENABLE_UPDATE"`
	// 当前这个消费者 配置的topic
	UpdateTopic string `toml:"update_topic" json:"update_topic" yaml:"update_topic"  env:"UPDATE_TOPIC"`
}

// func (s *TaskServiceImpl) AddAsyncTask(t *task.Task) {
// 	s.tasks.Add(t)
// }

// func (s *TaskServiceImpl) RemoveAsyncTask(t *task.Task) {
// 	news := types.New[*task.Task]()
// 	s.tasks.ForEach(func(existT *task.Task) {
// 		if t.Id != existT.Id {
// 			news.Add(t)
// 		}
// 	})
// 	s.tasks = news
// }

func (i *CronJobServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&cronjob.CronJob{})
		if err != nil {
			return err
		}
	}
	i.hostname, _ = os.Hostname()

	if i.EnableUpdate {
		i.updater = ioc_kafka.ConsumerGroup(i.hostname, []string{i.UpdateTopic})

		// 订阅更新事件
		go i.Run(i.ctx)
	}

	return nil
}

func (i *CronJobServiceImpl) Close(ctx context.Context) {
	i.ctx.Done()
}

func (i *CronJobServiceImpl) Name() string {
	return cronjob.APP_NAME
}
