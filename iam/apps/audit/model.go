package audit

import (
	"encoding/json"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/rs/xid"
)

func NewAuditLog() *AuditLog {
	return &AuditLog{
		Id:     xid.New().String(),
		Label:  map[string]string{},
		Extras: map[string]string{},
		Time:   time.Now().Unix(),
	}
}

// 用户操作事件
// 如何映射成 MongoDB BSON
type AuditLog struct {
	// 事件Id,
	Id string `json:"id" gorm:"primaryKey;column:id"`
	// 谁
	Who string `json:"who" gorm:"column:who"`
	// 在什么时间
	Time int64 `json:"time" gorm:"column:time"`
	// 操作人的Ip
	Ip string `json:"ip" gorm:"column:ip"`
	// User Agent
	UserAgent string `json:"user_agent" gorm:"column:user_agent"`

	// 做了什么操作,  服务:资源:动作
	// 服务 <cmdb, mcenter, ....>
	Service string `json:"service" gorm:"column:service"`
	// 哪个空间
	Namespace string `json:"namespace" gorm:"column:namespace"`
	// 资源 <secret, user, namespace, ...>
	ResourceType string `json:"resource_type" gorm:"column:resource_type"`
	// 动作 <list, get, update, create, delete, ....>
	Action string `json:"action" gorm:"column:action"`

	// 详情信息
	ResourceId string `json:"resource_id" gorm:"column:resource_id"`
	// 状态码 404
	StatusCode int `json:"status_code" gorm:"column:status_code"`
	// 具体信息
	ErrorMessage string `json:"error_message" gorm:"column:error_message"`

	// 标签
	Label map[string]string `json:"label" gorm:"column:label;serializer:json;type:json" description:"事件标签" optional:"true"`
	// 扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"额外属性" optional:"true"`
}

func (e *AuditLog) Load(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *AuditLog) TableName() string {
	return "audit_logs"
}

func (e *AuditLog) String() string {
	return pretty.ToJSON(e)
}

func (e *AuditLog) ToBusEvent(topic string) *bus.Event {
	data, _ := json.Marshal(e)
	return &bus.Event{
		Subject: topic,
		Data:    data,
	}
}
