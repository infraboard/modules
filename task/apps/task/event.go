package task

import (
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/task/apps/event"
)

func NewErrorEvent(msg string, taskId string) *types.Set[*event.EventSpec] {
	return types.New[*event.EventSpec]().Add(
		event.NewEventSpec().SetLevel(event.LEVEL_ERROR).SetMessage(msg).SetLabel(EVENT_LABLE_TASK_ID, taskId))
}

func NewInfoEvent(msg string, taskId string) *types.Set[*event.EventSpec] {
	return types.New[*event.EventSpec]().Add(
		event.NewEventSpec().SetLevel(event.LEVEL_INFO).SetMessage(msg).SetLabel(EVENT_LABLE_TASK_ID, taskId))
}
