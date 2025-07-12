package event

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	APP_NAME = "task_events"
)

func GetService() Service {
	return ioc.Controller().Get(APP_NAME).(Service)
}

type Service interface {
	// 添加事件
	AddEvent(context.Context, *EventSpec) (*Event, error)
	// 查询事件
	QueryEvent(context.Context, *QueryEventRequest) (*types.Set[*Event], error)
}

func NewQueryEventRequest() *QueryEventRequest {
	return &QueryEventRequest{
		PageRequest: *request.NewDefaultPageRequest(),
	}
}

type QueryEventRequest struct {
	// 分页参数
	request.PageRequest
	// 事件标签, TaskId
	Label map[string]string `json:"label" bson:"label" description:"事件标签" optional:"true"`
	// 排序方式
	OrderBy ORDER_BY `json:"order_by" bson:"order_by" description:"排序方式" optional:"true"`
}

func (r *QueryEventRequest) SetLabel(key, value string) *QueryEventRequest {
	if r.Label == nil {
		r.Label = map[string]string{}
	}
	r.Label[key] = value
	return r
}

func (r *QueryEventRequest) SetOrderBy(orerBy ORDER_BY) *QueryEventRequest {
	r.OrderBy = orerBy
	return r
}

type AddEventRequest struct {
}
