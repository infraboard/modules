package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/apps/webhook"
)

func (s *TaskServiceImpl) saveTask(ctx context.Context, ins *task.Task) {
	err := datasource.DBFromCtx(ctx).Save(ins).Error
	if err != nil {
		s.log.Error().Msgf("save task error, %s", err)
	}

	// 执行WebHook
	go s.runWebHook(ctx, ins)
}

func (s *TaskServiceImpl) updateTask(ctx context.Context, ins *task.Task) {
	err := datasource.DBFromCtx(ctx).Save(ins).Error
	if err != nil {
		s.log.Error().Msgf("save task error, %s", err)
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

func (s *TaskServiceImpl) saveEvent(ctx context.Context, events *types.Set[*event.EventSpec]) {
	_, err := event.GetService().AddEvent(ctx, events)
	if err != nil {
		s.log.Error().Msgf("save event error, %s", err)
	}
}
