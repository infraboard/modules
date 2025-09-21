package task

const (
	// 任务处于挂起, 等待排队
	STATUS_PENDDING STATUS = iota
	// 队列中
	STATUS_QUEUED
	// 任务正在运行
	STATUS_RUNNING
	// 任务正在更新
	STATUS_UPDATING
	// 取消中
	STATUS_CANCELING
	// 任务已完成
	STATUS_SUCCESS
	// 任务失败
	STATUS_FAILED
	// 任务已取消
	STATUS_CANCELED
	// 忽略执行
	STATUS_SKIPPED
)

var (
	STATUS_MAP = map[STATUS]string{
		STATUS_PENDDING: "PENDDING",
		STATUS_QUEUED:   "QUEUED",
		STATUS_RUNNING:  "RUNNING",
		STATUS_UPDATING: "UPDATING",
		STATUS_SUCCESS:  "SUCCESS",
		STATUS_FAILED:   "FAILED",
		STATUS_CANCELED: "CANCELED",
	}
	STATUS_COMPLETE = []STATUS{
		STATUS_SUCCESS,
		STATUS_FAILED,
		STATUS_CANCELED,
	}
)

func StatusCompleteString() []string {
	status := []string{}
	for _, s := range STATUS_COMPLETE {
		status = append(status, s.String())
	}
	return status
}

type STATUS int

func (s STATUS) String() string {
	return STATUS_MAP[s]
}

const (
	// 任务运行事件
	QUEUE_EVENT_TYPE_RUN QUEUE_EVENT_TYPE = "run"
	// 任务更新事件
	QUEUE_EVENT_TYPE_UPDATE QUEUE_EVENT_TYPE = "update"
	// 任务取消事件
	QUEUE_EVENT_TYPE_CANCEL QUEUE_EVENT_TYPE = "cancel"
)

type QUEUE_EVENT_TYPE string
