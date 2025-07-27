package impl

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

func (s *TaskServiceImpl) updateTaskStatus(ctx context.Context, ins *task.Task) {
	ins.SetUpdateAt(time.Now())
	err := datasource.DBFromCtx(ctx).Where("id = ?", ins.Id).Updates(&ins.TaskStatus).Error
	if err != nil {
		s.log.Error().Msgf("save task error, %s", err)
		return
	}

	// 执行WebHook
	go s.runWebHook(ctx, ins)
}

// 执行WebHook
func (s *TaskServiceImpl) runWebHook(ctx context.Context, ins *task.Task) {
	for _, hook := range ins.WebHooks {
		hook.RefTaskId = ins.Id

		// 判断是否需要Run
		if !hook.IsCondtionOk(ins.Status.String()) {
			s.log.Debug().Msgf("hook %s condition no ok", hook.TargetURL)
			continue
		}

		// 运行触发WebHook, 异步执行
		_, err := webhook.GetService().Run(ctx, &hook.WebHookSpec)
		if err != nil {
			s.log.Error().Msgf("hook %s run error, %s", hook.TargetURL, err)
			continue
		}
	}
}
