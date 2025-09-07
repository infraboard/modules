package task

import (
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/task/apps/event"
)

func NewDebugEvent(taskId string, format string, a ...any) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_DEBUG).SetMessage(fmt.Sprintf(format, a...)).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewInfoEvent(taskId string, format string, a ...any) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_INFO).SetMessage(fmt.Sprintf(format, a...)).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewWarnEvent(taskId string, format string, a ...any) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_WARN).SetMessage(fmt.Sprintf(format, a...)).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewErrorEvent(taskId string, format string, a ...any) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_ERROR).SetMessage(fmt.Sprintf(format, a...)).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewQueueEvent() *QueueEvent {
	return &QueueEvent{
		Type: QUEUE_EVENT_TYPE_RUN,
	}
}

type QueueEvent struct {
	// 事件类型
	Type QUEUE_EVENT_TYPE `json:"type"`
	// 事件值
	TaskId string `json:"task_id"`
}

func (e *QueueEvent) Load(v []byte) error {
	return json.Unmarshal(v, e)
}

func (e *QueueEvent) String() string {
	return pretty.ToJSON(e)
}
