package event

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	AddEvent(context.Context, *types.Set[*EventSpec]) (*types.Set[*Event], error)
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
}

type AddEventRequest struct {
}

func NewEvent(spec EventSpec) *Event {
	return &Event{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		EventSpec: spec,
	}
}

type Event struct {
	// 任务Id
	Id string `json:"id" gorm:"column:id;type:uint;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// 事件定义
	EventSpec
}

func (e *Event) TableName() string {
	return "task_events"
}

func NewEventSpec() *EventSpec {
	return &EventSpec{
		Time:  time.Now(),
		Level: LEVEL_DEBUG,
		Label: map[string]string{},
	}
}

type EventSpec struct {
	// 事件发生时间
	Time time.Time `json:"time" gorm:"column:time;type:timestamp;" description:"事件发生时间"`
	// 事件的级别
	Level LEVEL `json:"level" gorm:"column:level;type:tinyint(2);" description:"事件的级别"`
	// 事件信息
	Message string `json:"message" gorm:"column:message;type:text;" description:"事件信息"`
	// 事件详情
	Detail string `json:"detail" gorm:"column:detail;type:text;" description:"事件详情"`
	// 事件标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json;type:json" description:"事件标签" optional:"true"`
}

func (s *EventSpec) SetLevel(l LEVEL) *EventSpec {
	s.Level = l
	return s
}

func (s *EventSpec) SetMessage(msg string) *EventSpec {
	s.Message = msg
	return s
}

func (s *EventSpec) SetLabel(key, value string) *EventSpec {
	if s.Label == nil {
		s.Label = map[string]string{}
	}
	s.Label[key] = value
	return s
}
