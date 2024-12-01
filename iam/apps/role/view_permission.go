package role

import (
	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/view"
)

func NewViewPermission(roleId uint64, spec *ViewPermissionSpec) *ViewPermission {
	return &ViewPermission{
		ResourceMeta:       *apps.NewResourceMeta(),
		RoleId:             roleId,
		ViewPermissionSpec: *spec,
	}
}

type ViewPermission struct {
	// 基础数据
	apps.ResourceMeta
	// 角色Id
	RoleId uint64 `json:"role_id" gorm:"column:role_id;index" description:"Role Id"`
	// Menu权限定义
	ViewPermissionSpec
}

func (r *ViewPermission) TableName() string {
	return "view_permissions"
}

type ViewPermissionSpec struct {
	// 创建者ID
	CreateBy uint64 `json:"create_by" gorm:"column:create_by" description:"创建者ID" optional:"true"`
	// 角色描述
	Description string `json:"description" gorm:"column:description;type:text" bson:"description" description:"角色描述"`
	// 服务
	Service string `json:"service" gorm:"column:service;type:varchar(100);index" bson:"service" description:"服务名称"`
	// 页面标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index" description:"页面标签"`
	// 视图类型
	ViewType view.TYPE `json:"view_type" gorm:"column:view_type;type:tinyint(1);index" description:"视图类型"`
	// 视图路径, 如果是Menu就是Menu路径，如果是Page就是Page路径
	ViewPath string `json:"view_path" gorm:"column:view_path;type:varchar(200);index" bson:"view_path" description:"视图路径(可以通配)"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息" optional:"true"`
}
