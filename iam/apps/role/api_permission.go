package role

import (
	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/endpoint"
)

type ApiPermission struct {
	// 基础数据
	apps.Meta
	// Api权限定义
	ApiPermissionSpec
}

func (r *ApiPermission) TableName() string {
	return "api_permissions"
}

type ApiPermissionSpec struct {
	// 角色Id
	RoleId uint64 `json:"role_id" gorm:"column:role_id;index"`
	// 权限序号
	SequenceNumber uint16 `json:"sequence_number" gorm:"column:sequence_number" description:"权限序号" optional:"true"`
	// 创建者ID
	CreateBy uint64 `json:"create_by" gorm:"column:create_by" description:"创建者ID" optional:"true"`
	// 角色描述
	Description string `json:"description" gorm:"column:description;type:text" bson:"description" description:"角色描述"`
	// 效力
	Effect EFFECT_TYPE `json:"effect" gorm:"column:effect;type:tinyint(1);index" bson:"effect" description:"效力"`
	// 服务
	Service string `json:"service" gorm:"column:service;type:varchar(100);index" bson:"service" description:"服务名称"`
	// 资源列表
	Resource string `json:"resource" gorm:"column:resource;type:varchar(100);index" bson:"resource" description:"资源名称"`
	// 资源操作
	Action string `json:"action" bson:"action" gorm:"column:action;type:varchar(100);index"`
	// 读或者写
	ActionType endpoint.ACTION_TYPE `json:"action_type" bson:"action_type" gorm:"column:action_type;type:tinyint(1);index"`
	// 操作标签
	ActionLabel string `json:"action_label" gorm:"column:action_label;type:varchar(200);index"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息" optional:"true"`
}
