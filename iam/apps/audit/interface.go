package audit

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

var (
	AppName = "audits"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 存储
	SaveAuditLog(context.Context, *AuditLog) error
	// 查询
	QueryAuditLog(context.Context, *QueryAuditLogRequest) (*types.Set[*AuditLog], error)
}

func NewQueryAuditLogRequest() *QueryAuditLogRequest {
	return &QueryAuditLogRequest{
		PageRequest: request.NewDefaultPageRequest(),
		Label:       map[string]string{},
	}
}

type QueryAuditLogRequest struct {
	// 分页请求参数
	*request.PageRequest
	// 事件标签, TaskId
	Label map[string]string `json:"label" bson:"label" description:"事件标签" optional:"true"`
}
