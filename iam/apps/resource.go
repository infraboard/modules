package apps

import (
	"time"
)

func NewMeta() *Meta {
	return &Meta{
		CreatedAt: time.Now(),
	}
}

type Meta struct {
	// 自增Id
	Id uint64 `json:"id" gorm:"column:id;type:uint;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;" optional:"true" description:"更新时间"`
	// 删除时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index" optional:"true" description:"删除时间"`
}
