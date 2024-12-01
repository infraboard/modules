package endpoint

import (
	"fmt"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/modules/iam/apps"
)

func NewEndpoint() *Endpoint {
	return &Endpoint{
		ResourceMeta: *apps.NewResourceMeta(),
	}
}

// Endpoint Service's features
type Endpoint struct {
	// 基础数据
	apps.ResourceMeta
	// 路由条目信息
	RouteEntry `json:"route_entry" bson:",inline" validate:"required"`
}

func (u *Endpoint) TableName() string {
	return "namespaces"
}

func (u *Endpoint) SetRouteEntry(v RouteEntry) *Endpoint {
	u.RouteEntry = v
	return u
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
	Resource string `json:"resource" bson:"resource" gorm:"column:resource;type:varchar(100);index"`
	// 资源操作
	Action string `json:"action" bson:"action" gorm:"column:action;type:varchar(100);index"`
	// 读或者写
	AccessMode ACCESS_MODE `json:"access_mode" bson:"access_mode" gorm:"column:access_mode;type:tinyint(1);index"`
	// 操作标签
	ActionLabel string `json:"action_label" gorm:"column:action_label;type:varchar(200);index"`
	// 函数名称
	FunctionName string `json:"function_name" bson:"function_name" gorm:"column:function_name;type:varchar(100)"`
	// HTTP path 用于自动生成http api
	Path string `json:"path" bson:"path" gorm:"column:path;type:varchar(200);index"`
	// HTTP method 用于自动生成http api
	Method string `json:"method" bson:"method" gorm:"column:method;type:varchar(100);index"`
	// 接口说明
	Description string `json:"description" bson:"description" gorm:"column:description;type:text"`
	// 是否校验用户身份 (acccess_token 校验)
	RequiredAuth bool `json:"required_auth" bson:"required_auth" gorm:"column:required_auth;type:tinyint(1)"`
	// 验证码校验(开启双因子认证需要) (code 校验)
	RequiredCode bool `json:"required_code" bson:"required_code" gorm:"column:required_code;type:tinyint(1)"`
	// 开启鉴权
	RequiredPerm bool `json:"required_perm" bson:"required_perm" gorm:"column:required_perm;type:tinyint(1)"`
	// ACL模式下, 允许的通过的身份标识符, 比如角色, 用户类型之类
	RequiredRole []string `json:"required_role" bson:"required_role" gorm:"column:required_role;serializer:json;type:json"`
	// 是否开启操作审计, 开启后这次操作将被记录
	RequiredAudit bool `json:"required_audit" bson:"required_audit" gorm:"column:required_audit;type:tinyint(1)"`
	// 名称空间不能为空
	RequiredNamespace bool `json:"required_namespace" bson:"required_namespace" gorm:"column:required_namespace;type:tinyint(1)"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json;type:json"`
}

func (e *RouteEntry) LoadMeta(meta map[string]any) {
	e.Resource = GetRouteMeta[string](meta, META_RESOURCE_KEY)
	e.RequiredAuth = GetRouteMeta[bool](meta, META_REQUIRED_AUTH_KEY)
	e.RequiredCode = GetRouteMeta[bool](meta, META_REQUIRED_CODE_KEY)
	e.RequiredPerm = GetRouteMeta[bool](meta, META_REQUIRED_PERM_KEY)
	e.RequiredRole = GetRouteMeta[[]string](meta, META_REQUIRED_ROLE_KEY)
	e.RequiredAudit = GetRouteMeta[bool](meta, META_REQUIRED_AUDIT_KEY)
	e.RequiredNamespace = GetRouteMeta[bool](meta, META_REQUIRED_NAMESPACE_KEY)
}

// UniquePath todo
func (e *RouteEntry) HasRequiredRole() bool {
	return len(e.RequiredRole) > 0
}

// UniquePath todo
func (e *RouteEntry) UniquePath() string {
	return fmt.Sprintf("%s.%s", e.Method, e.Path)
}

func (e *RouteEntry) IsRequireRole(target string) bool {
	for i := range e.RequiredRole {
		if e.RequiredRole[i] == "*" {
			return true
		}

		if e.RequiredRole[i] == target {
			return true
		}
	}

	return false
}

func (e *RouteEntry) SetRequiredAuth(v bool) *RouteEntry {
	e.RequiredAuth = v
	return e
}

func (e *RouteEntry) AddRequiredRole(roles ...string) *RouteEntry {
	e.RequiredRole = append(e.RequiredRole, roles...)
	return e
}

func (e *RouteEntry) SetRequiredPerm(v bool) *RouteEntry {
	e.RequiredPerm = v
	return e
}

func (e *RouteEntry) SetLabel(value string) *RouteEntry {
	e.ActionLabel = value
	return e
}

func (e *RouteEntry) SetExtensionFromMap(m map[string]string) *RouteEntry {
	if e.Extras == nil {
		e.Extras = map[string]string{}
	}

	for k, v := range m {
		e.Extras[k] = v
	}
	return e
}

func (e *RouteEntry) SetRequiredCode(v bool) *RouteEntry {
	e.RequiredCode = v
	return e
}

func NewEntryFromRestRequest(req *restful.Request) *RouteEntry {
	entry := NewRouteEntry()

	// 请求拦截
	route := req.SelectedRoute()
	if route == nil {
		return nil
	}

	entry.FunctionName = route.Operation()
	entry.Method = route.Method()
	entry.LoadMeta(route.Metadata())
	entry.Path = route.Path()
	return entry
}

func NewEntryFromRestRoute(route restful.RouteReader) *RouteEntry {
	entry := NewRouteEntry()
	entry.FunctionName = route.Operation()
	entry.Method = route.Method()
	entry.LoadMeta(route.Metadata())
	entry.Path = route.Path()

	entry.Path = entry.UniquePath()
	return entry
}
