package task

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	APP_NAME = "tasks"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

type Service interface {
	// 任务执行
	Run(context.Context, *TaskSpec) *Task
	// 查询任务列表
	QueryTask(context.Context, *QueryTaskRequest) (*types.Set[*Task], error)
	// 查询任务详情
	DescribeTask(context.Context, *DescribeTaskRequest) (*Task, error)
}

func NewDescribeTaskRequest(taskId string) *DescribeTaskRequest {
	return &DescribeTaskRequest{
		TaskId: taskId,
	}
}

type DescribeTaskRequest struct {
	TaskId string `json:"task_id"`
}

type QueryTaskRequest struct {
}
