package cronjob

import (
	"encoding/json"

	"github.com/infraboard/mcube/v2/tools/pretty"
)

func NewQueueEvent() *QueueEvent {
	return &QueueEvent{}
}

type QueueEvent struct {
	// 事件类型
	Type QUEUE_EVENT_TYPE `json:"type"`
	// 事件值
	CronJobId string `json:"cron_job_id"`
}

func (e *QueueEvent) Load(v []byte) error {
	return json.Unmarshal(v, e)
}

func (e *QueueEvent) String() string {
	return pretty.ToJSON(e)
}
