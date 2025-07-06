package task

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

var (
	STATUS_MAP = map[STATUS]string{
		STATUS_PENDING:  "PENDDING",
		STATUS_RUNNING:  "RUNNING",
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
	// 任务是一个函数
	TYPE_FUNCTION TYPE = "function"
)

type TYPE string
