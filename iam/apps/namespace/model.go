package namespace

import "github.com/infraboard/modules/iam/apps"

func NewNmespace() *Nmespace {
	return &Nmespace{
		Meta: *apps.NewMeta().WithUUID(),
	}
}

type Nmespace struct {
	// 基础数据
	apps.Meta
	// 空间属性
	CreateNamespaceRequest
}

func (u *Nmespace) TableName() string {
	return "namespaces"
}

func NewCreateNamespaceRequest() *CreateNamespaceRequest {
	return &CreateNamespaceRequest{
		Extras: map[string]string{},
	}
}

type CreateNamespaceRequest struct {
	// 父Namespace Id
	ParentId uint64 `json:"parent_id" bson:"parent_id" gorm:"column:parent_id;type:uint;index"`
	// 全局唯一
	Name string `json:"name" bson:"name" validate:"required" gorm:"column:name;type:varchar(200);not null;uniqueIndex"`
	// 空间负责人
	Owner uint64 `json:"owner" bson:"owner" gorm:"column:owner;type:uint;index;not null"`
	// 禁用项目, 该项目所有人暂时都无法访问
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled;type:tinyint(1)"`
	// 项目描述图片
	Picture string `json:"picture" bson:"picture" gorm:"column:picture;type:varchar(200)"`
	// 项目描述
	Description string `json:"description" bson:"description" gorm:"column:description;type:text"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json;type:json"`
}
