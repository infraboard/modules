package impl

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&TaskServiceImpl{
		ctx:         context.Background(),
		tasks:       types.New[*task.Task](),
		RunTopic:    "task_run_events",
		CancelTopic: "task_cancel_events",
	})
}

var _ task.Service = (*TaskServiceImpl)(nil)

type TaskServiceImpl struct {
	ioc.ObjectImpl

	// 日志
	log *zerolog.Logger
	// 允许时上下文
	ctx context.Context
	// 异步任务的状态, 运行中的
	tasks *types.Set[*task.Task]
	// 节点名称
	node_name string

	// 随机节点名, 用于单节点调试
	RandomNodeName bool `toml:"random_node_name" json:"random_node_name" yaml:"random_node_name"  env:"RANDOM_NODE_NAME"`
	// 运行队列名称
	RunTopic string `toml:"run_topic" json:"run_topic" yaml:"run_topic"  env:"RUN_TOPIC"`
	// 取消队列名称
	CancelTopic string `toml:"cancel_topic" json:"cancel_topic" yaml:"cancel_topic"  env:"CANCEL_TOPIC"`
}

func (s *TaskServiceImpl) AddAsyncTask(t *task.Task) {
	s.tasks.Add(t)
}

func (s *TaskServiceImpl) GetAsyncTask(id string) *task.Task {
	for _, item := range s.tasks.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (s *TaskServiceImpl) RemoveAsyncTask(t *task.Task) {
	news := types.New[*task.Task]()
	s.tasks.ForEach(func(existT *task.Task) {
		if t.Id != existT.Id {
			news.Add(t)
		}
	})
	s.tasks = news
}

func (i *TaskServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&task.Task{})
		if err != nil {
			return err
		}
	}

	// 设置节点名称
	if i.RandomNodeName {
		i.node_name = uuid.NewString()
	} else {
		i.node_name, _ = os.Hostname()
	}

	// 处理取消事件
	i.log.Info().Msgf("cancel topic: %s", i.CancelTopic)
	go i.HandleCancelEvents(i.ctx)

	// 处理运行事件
	i.log.Info().Msgf("run topic: %s", i.RunTopic)
	go i.HandleRunEvents(i.ctx)
	return nil
}

func (i *TaskServiceImpl) Name() string {
	return task.APP_NAME
}
