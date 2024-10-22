package entrypoint

import "github.com/infraboard/modules/iam/apps"

func NewEndpoint() *Endpoint {
	return &Endpoint{
		Meta: *apps.NewMeta(),
	}
}

// Endpoint Service's features
type Endpoint struct {
	// 基础数据
	apps.Meta
	// 路由条目信息
	RouteEntry `json:"route_entry" bson:",inline" validate:"required"`
}

func (u *Endpoint) TableName() string {
	return "namespaces"
}

func NewRouteEntry() *RouteEntry {
	return &RouteEntry{
		RequiredRole: []string{},
		Extras:       map[string]string{},
	}
}

// Entry 路由条目
type RouteEntry struct {
	// 该功能属于那个服务
	Service string `json:"service" bson:"service" validate:"required,lte=64" gorm:"column:service;type:varchar(100);index"`
	// 服务那个版本的功能
	Version string `json:"version" bson:"version" validate:"required,lte=64" gorm:"column:version;type:varchar(100)"`
	// 资源名称
	Resource string `json:"resource" bson:"resource" gorm:"column:resource;type:varchar(100)"`
	// 函数名称
	FunctionName string `json:"function_name" bson:"function_name" gorm:"column:function_name;type:varchar(100)"`
	// HTTP path 用于自动生成http api
	Path string `json:"path" bson:"path" gorm:"column:path;type:varchar(200);index"`
	// HTTP method 用于自动生成http api
	Method string `json:"method" bson:"method" gorm:"column:method;type:varchar(100);index"`
	// 接口说明
	Description string `json:"description" bson:"description" gorm:"column:description;type:text"`
	// 是否校验用户身份 (acccess_token 校验)
	AuthEnable bool `json:"auth_enable" bson:"auth_enable" gorm:"column:auth_enable;type:tinyint(1)"`
	// 验证码校验(开启双因子认证需要) (code 校验)
	CodeEnable bool `json:"code_enable" bson:"code_enable" gorm:"column:code_enable;type:tinyint(1)"`
	// 开启鉴权
	PermEnable bool `json:"perm_enable" bson:"perm_enable" gorm:"column:perm_enable;type:tinyint(1)"`
	// ACL模式下, 允许的通过的身份标识符, 比如角色, 用户类型之类
	RequiredRole []string `json:"required_role" bson:"required_role" gorm:"column:required_role;serializer:json;type:json"`
	// 是否开启操作审计, 开启后这次操作将被记录
	AuditLog bool `json:"audit_log" bson:"audit_log" gorm:"column:audit_log;type:tinyint(1)"`
	// 名称空间不能为空
	RequiredNamespace bool `json:"required_namespace" bson:"required_namespace" gorm:"column:required_namespace;type:tinyint(1)"`
	// 标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json;type:json"`
}
