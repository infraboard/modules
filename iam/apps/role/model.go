package role

import (
	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/menu"
)

func NewRole() *Role {
	return &Role{
		Meta:      *apps.NewMeta(),
		Menus:     []*menu.Menu{},
		Endpoints: []*endpoint.Endpoint{},
	}
}

type Role struct {
	// 基础数据
	apps.Meta
	// 角色创建信息
	CreateRoleRequest
	// 菜单
	Menus []*menu.Menu `json:"menus,omitempty" gorm:"-" description:"角色关联的菜单"`
	// API
	Endpoints []*endpoint.Endpoint `json:"endpoints,omitempty" gorm:"-" description:"角色关联的API"`
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
	CreateBy uint64 `json:"create_by" gorm:"column:create_by" description:"创建者ID" optional:"true"`
	// 角色名称
	Name string `json:"name" gorm:"column:name;type:varchar(100);index" bson:"name" description:"角色名称"`
	// 角色描述
	Description string `json:"description" gorm:"column:description;type:text" bson:"description" description:"角色描述"`
	// 是否启用
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled;type:tinyint(1)" description:"是否启用" optional:"true"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index" description:"标签" optional:"true"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息" optional:"true"`
}

type RoleAssociateMenuRecord struct {
	// 基础数据
	apps.Meta
	// 角色Id
	RoleId uint64 `json:"role_id" gorm:"column:role_id;index"`
	// 菜单Id
	MenuId uint64 `json:"menu_id" gorm:"column:menu_id;index"`
	// 关联Menu
	Menu *menu.Menu `json:"menu" gorm:"-"`
}

type RoleAssociateEndpointRecord struct {
	// 基础数据
	apps.Meta
	// 角色Id
	RoleId uint64 `json:"role_id" gorm:"column:role_id;index"`
	// 菜单Id
	EndpointId uint64 `json:"endpoint_id" gorm:"column:endpoint_id;index"`
	// 关联API
	Endpoint *endpoint.Endpoint `json:"endpoint" gorm:"-"`
}
