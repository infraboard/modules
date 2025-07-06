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
	// 执行WebHook
	Run(context.Context, *WebHookSpec) *WebHook
	// 查询WebHook执行记录
	QueryWebHook(context.Context, *QueryWebHookRequest) (*types.Set[*WebHook], error)
}

type QueryWebHookRequest struct {
	// 关联Task
	RefTaskId string `json:"ref_task_id"`
}
