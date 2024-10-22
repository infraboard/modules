package role

import "github.com/infraboard/modules/iam/apps"

func NewRole() *Role {
	return &Role{
		Meta: *apps.NewMeta(),
	}
}

type Role struct {
	// 基础数据
	apps.Meta
	// 角色创建信息
	CreateRoleRequest
}

func (u *Role) TableName() string {
	return "roles"
}

func NewCreateRoleRequest() *CreateRoleRequest {
	return &CreateRoleRequest{
		Extras: map[string]string{},
	}
}

type CreateRoleRequest struct {
	// 创建者ID
	CreateBy uint64 `json:"create_by" gorm:"column:create_by"`
	// 角色名称
	Name string `json:"name" gorm:"column:name;type:varchar(100);index" bson:"name"`
	// 角色描述
	Description string `json:"description" gorm:"column:description;type:text" bson:"description"`
	// 是否启用
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled;type:tinyint(1)"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json"`
}
