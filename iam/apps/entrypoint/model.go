package entrypoint

// Endpoint Service's features
type Endpoint struct {
	// 端点名称
	Id string `json:"id" bson:"_id" validate:"required,lte=64" gorm:"column:id"`
	// 创建时间
	CreateAt int64 `json:"create_at" bson:"create_at" gorm:"column:create_at"`
	// 更新时间
	UpdateAt int64 `json:"update_at" bson:"update_at" gorm:"column:update_at"`
	// 该功能属于那个服务
	Service string `json:"service" bson:"service" validate:"required,lte=64" gorm:"column:service"`
	// 服务那个版本的功能
	Version string `json:"version" bson:"version" validate:"required,lte=64" gorm:"column:version"`
	// 路由条目信息
	*RouteEntry `json:"route_entry" bson:",inline" validate:"required"`
}

// Entry 路由条目
type RouteEntry struct {
	// 函数名称
	FunctionName string `json:"function_name" bson:"function_name" gorm:"column:function_name"`
	// HTTP path 用于自动生成http api
	Path string `json:"path" bson:"path" gorm:"column:path"`
	// HTTP method 用于自动生成http api
	Method string `json:"method" bson:"method" gorm:"column:method"`
	// 资源名称
	Resource string `json:"resource" bson:"resource" gorm:"column:resource"`
	// 是否校验用户身份 (acccess_token 校验)
	AuthEnable bool `json:"auth_enable" bson:"auth_enable" gorm:"column:auth_enable"`
	// 验证码校验(开启双因子认证需要) (code 校验)
	CodeEnable bool `json:"code_enable" bson:"code_enable" gorm:"column:code_enable"`
	// 权限验证的模式, 支持ACL/PRBAC, 默认PRBAC
	PermissionMode string `json:"permission_mode" bson:"permission_mode" gorm:"column:permission_mode"`
	// PRBAC模式下 开启权限校验
	PermissionEnable bool `json:"permission_enable" bson:"permission_enable" gorm:"column:permission_enable"`
	// ACL模式下, 允许的通过的身份标识符, 比如角色, 用户类型之类
	RequiredRole []string `json:"required_role" bson:"required_role" gorm:"column:required_role"`
	// 是否开启操作审计, 开启后这次操作将被记录
	AuditLog bool `json:"audit_log" bson:"audit_log" gorm:"column:audit_log"`
	// 名称空间不能为空
	RequiredNamespace bool `json:"required_namespace" bson:"required_namespace" gorm:"column:required_namespace"`
	// 标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json"`
	// 扩展属性
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json"`
}
