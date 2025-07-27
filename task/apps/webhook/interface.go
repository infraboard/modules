package webhook

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	APP_NAME = "webhooks"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

type Service interface {
	// 执行WebHook, 触发执行, 执行状态通过HookId查询
	Run(context.Context, *WebHookSpec) (*WebHook, error)
	// 查询WebHook具体执行状态
	DescribeWebHook(context.Context, *DescribeWebHookRequest) (*WebHook, error)
	// 查询WebHook执行记录列表
	QueryWebHook(context.Context, *QueryWebHookRequest) (*types.Set[*WebHook], error)
}

type DescribeWebHookRequest struct {
	Id string
}

func NewQueryWebHookRequest() *QueryWebHookRequest {
	return &QueryWebHookRequest{}
}

type QueryWebHookRequest struct {
	// 关联Task
	RefTaskId string `json:"ref_task_id"`
}

func (r *QueryWebHookRequest) SetRefTaskId(taskId string) *QueryWebHookRequest {
	r.RefTaskId = taskId
	return r
}
