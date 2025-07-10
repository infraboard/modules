package impl

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&CronServiceImpl{
		// tasks: types.New[*task.Task](),
	})
}

var _ cronjob.Service = (*CronServiceImpl)(nil)

type CronServiceImpl struct {
	ioc.ObjectImpl

	log *zerolog.Logger

	// 异步任务的状态, 运行中的
	// tasks *types.Set[*task.Task]
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

func (i *CronServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&cronjob.CronJob{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *CronServiceImpl) Name() string {
	return cronjob.APP_NAME
}
