package apps

import (
	"time"

	"github.com/infraboard/mcube/v2/types/uuid"
)

func NewMeta() *Meta {
	return &Meta{
		CreatedAt: time.Now(),
	}
}

type Meta struct {
	// 自增Id
	Id uint64 `json:"id" gorm:"column:id;type:uint;primary_key;" unique:"true"`
	// UUID
	UUID *uuid.BinaryUUID `type:"string" swagger:"required" json:"uuid" bson:"uuid" gorm:"column:uuid;type:binary(16);uniqueIndex" unique:"true"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;" optional:"true"`
	// 删除时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index" optional:"true"`
}

func (m *Meta) WithUUID() *Meta {
	v := uuid.NewUUID()
	m.UUID = &v
	return m
}
