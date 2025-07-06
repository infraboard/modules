package webhook

const (
	// 任务处于挂起
	STATUS_PENDING STATUS = iota
	// 任务正在运行
	STATUS_RUNNING
	// 任务已完成
	STATUS_SUCCESS
	// 任务失败
	STATUS_FAILED
	// 任务已取消
	STATUS_CANCELED
)

type STATUS int
