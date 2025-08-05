package impl

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&CronJobServiceImpl{
		EnableUpdate:   true,
		RandomNodeName: false,
		UpdateTopic:    "cronjob_update_events",
		TaskLockTTL:    "5m",
		ctx:            context.Background(),
	})
}

var _ cronjob.Service = (*CronJobServiceImpl)(nil)

type CronJobServiceImpl struct {
	ioc.ObjectImpl

	// 日志
	log *zerolog.Logger
	// 节点名称
	node_name string

	// 允许时上下文
	ctx context.Context

	// 随机节点名, 用于单节点调试
	RandomNodeName bool `toml:"random_node_name" json:"random_node_name" yaml:"random_node_name"  env:"RANDOM_NODE_NAME"`
	// 启用更新队列
	EnableUpdate bool `toml:"enable_update" json:"enable_update" yaml:"enable_update"  env:"ENABLE_UPDATE"`
	// 当前这个消费者 配置的topic
	UpdateTopic string `toml:"update_topic" json:"update_topic" yaml:"update_topic"  env:"UPDATE_TOPIC"`
	// 执行时锁的TTL
	TaskLockTTL string `toml:"task_lock_ttl" json:"task_lock_ttl" yaml:"task_lock_ttl"  env:"TASK_LOCK_TTL"`

	taskLockTTL time.Duration
}

func (i *CronJobServiceImpl) Priority() int {
	return cronjob.PRIORITY
}

func (i *CronJobServiceImpl) Init() error {

	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&cronjob.CronJob{})
		if err != nil {
			return err
		}
	}

	lttl, err := time.ParseDuration(i.TaskLockTTL)
	if err != nil {
		return err
	}
	i.taskLockTTL = lttl

	// 设置节点名称
	if i.RandomNodeName {
		i.node_name = uuid.NewString()
	} else {
		i.node_name, _ = os.Hostname()
	}

	if i.EnableUpdate {
		// 订阅更新事件
		i.HandleUpdateEvents(i.ctx)
	}

	return nil
}

func (i *CronJobServiceImpl) Close(ctx context.Context) {
	i.ctx.Done()
}

func (i *CronJobServiceImpl) Name() string {
	return cronjob.APP_NAME
}
