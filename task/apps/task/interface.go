package task

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
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
	// 创建任务
	CreateTask(context.Context, *TaskSpec) (*Task, error)
	// 任务执行, 包含同步和异步
	Run(context.Context, *TaskSpec) (*Task, error)
	// 触发任务状态更新, 异步任务时需要使用
	UpdateTaskStatus(context.Context, *UpdateTaskStatusRequest) (*Task, error)
	// 任务取消
	Cancel(context.Context, *CancelRequest) (*Task, error)
	// 查询任务列表
	QueryTask(context.Context, *QueryTaskRequest) (*types.Set[*Task], error)
	// 查询任务详情
	DescribeTask(context.Context, *DescribeTaskRequest) (*Task, error)
}

func NewUpdateTaskStatusRequest(taskId string) *UpdateTaskStatusRequest {
	return &UpdateTaskStatusRequest{
		DescribeTaskRequest: *NewDescribeTaskRequest(taskId),
	}
}

type UpdateTaskStatusRequest struct {
	DescribeTaskRequest
}

type CancelRequest struct {
	DescribeTaskRequest
}

func NewDescribeTaskRequest(taskId string) *DescribeTaskRequest {
	return &DescribeTaskRequest{
		TaskId: taskId,
	}
}

type DescribeTaskRequest struct {
	TaskId string `json:"task_id"`
}

func NewQueryTaskRequest() *QueryTaskRequest {
	return &QueryTaskRequest{
		PageRequest: *request.NewDefaultPageRequest(),
		Label:       map[string]string{},
	}
}

type QueryTaskRequest struct {
	request.PageRequest

	Label map[string]string `json:"label"`
}

func (r *QueryTaskRequest) SetLabel(key, value string) *QueryTaskRequest {
	if r.Label == nil {
		r.Label = map[string]string{}
	}

	r.Label[key] = value
	return r
}
