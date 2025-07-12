package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/task/apps/event"
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

		// 运行WebHook
		wh := webhook.GetService().Run(ctx, &hook.WebHookSpec)
		hook.Status = wh.Status
		switch hook.Status {
		case webhook.STATUS_SUCCESS:
			s.saveEvent(ctx, task.NewInfoEvent(fmt.Sprintf("web hook %s exec success", hook.TargetURL), ins.Id))
		default:
			s.saveEvent(ctx, task.NewErrorEvent(fmt.Sprintf("web hook %s exec failed", hook.TargetURL), ins.Id))
		}
	}
}

func (s *TaskServiceImpl) saveEvent(ctx context.Context, e *event.EventSpec) {
	_, err := event.GetService().AddEvent(ctx, e)
	if err != nil {
		s.log.Error().Msgf("save event error, %s", err)
	}
}
