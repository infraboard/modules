package event

import (
	"time"

	"github.com/google/uuid"
)

func NewEvent(spec EventSpec) *Event {
	return &Event{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		EventSpec: spec,
	}
}

type Event struct {
	// 任务Id
	Id string `json:"id" gorm:"column:id;type:string;primary_key;" unique:"true" description:"Id"`
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
	// 事件所属资源
	Resource string `json:"resource" gorm:"column:resource;type:varchar(120);" description:"事件所属资源"`
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

func (s *EventSpec) SetDetail(detail string) *EventSpec {
	s.Detail = detail
	return s
}

func (s *EventSpec) SetLabel(key, value string) *EventSpec {
	if s.Label == nil {
		s.Label = map[string]string{}
	}
	s.Label[key] = value
	return s
}
