package apps

import (
	"strconv"
	"time"
)

func NewResourceMeta() *ResourceMeta {
	return &ResourceMeta{
		CreatedAt: time.Now(),
	}
}

type ResourceMeta struct {
	// 自增Id
	Id uint64 `json:"id" gorm:"column:id;type:uint;primary_key;" unique:"true" description:"Id"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;" description:"创建时间"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;" optional:"true" description:"更新时间"`
	// 删除时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index" optional:"true" description:"删除时间"`
}

type GetRequest struct {
	Id uint64 `json:"id"`
}

func (r *GetRequest) SetId(id uint64) {
	r.Id = id
}

func (r *GetRequest) SetIdByString(id string) error {
	v, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	r.Id = v
	return nil
}

type Scope struct {
	// 空间
	NamespaceId *uint64 `json:"namespace_id" bson:"namespace_id" gorm:"column:namespace_id;type:varchar(200);index" description:"空间Id" optional:"true"`
	// 访问范围定义, 比如 Env=Dev
	Scope map[string]string `json:"scope" bson:"scope" gorm:"column:scope;serializer:json;type:json" description:"数据访问的范围定义" optional:"true"`
}
