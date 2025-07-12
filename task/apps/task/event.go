package task

import (
	"encoding/json"

	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/task/apps/event"
)

func NewErrorEvent(msg string, taskId string) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_ERROR).SetMessage(msg).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewInfoEvent(msg string, taskId string) *event.EventSpec {
	return event.NewEventSpec().SetLevel(event.LEVEL_INFO).SetMessage(msg).SetLabel(EVENT_LABLE_TASK_ID, taskId)
}

func NewQueueEvent() *QueueEvent {
	return &QueueEvent{}
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
