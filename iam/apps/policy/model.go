package policy

import "github.com/infraboard/mcube/v2/pb/resource"

type Policy struct {
	// 端点名称
	Id string `json:"id" bson:"_id" validate:"required,lte=64" gorm:"column:id"`
	// 创建时间
	CreateAt int64 `json:"create_at" bson:"create_at" gorm:"column:create_at"`
	// 更新时间
	UpdateAt int64 `json:"update_at" bson:"update_at" gorm:"column:update_at"`
	// 策略定义
	*CreatePolicyRequest
}

type CreatePolicyRequest struct {
	// 创建者
	CreateBy string `json:"create_by" bson:"create_by" gorm:"column:create_by"`
	// 范围
	Namespace string `json:"namespace" bson:"namespace" gorm:"column:namespace" validate:"lte=120"`
	// 用户Id
	UserId string `json:"user_id" bson:"user_id" gorm:"column:user_id" validate:"required,lte=120"`
	// 角色Id
	RoleId string `json:"role_id" bson:"role_id" gorm:"column:role_id" validate:"required,lte=40"`
	// 该角色的生效范围
	Scope []*resource.LabelRequirement `json:"scope" bson:"scope" gorm:"column:scope"`
	// 策略过期时间
	ExpiredTime int64 `json:"expired_time" bson:"expired_time" gorm:"column:expired_time"`
	// 只读策略, 不允许用户修改, 一般用于系统管理
	ReadOnly bool `json:"read_only" bson:"read_only" gorm:"column:read_only"`
	// 启用该策略
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled"`
	// 策略标签
	Label map[string]string `json:"label" bson:"label" gorm:"column:label;serializer:json"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json"`
}
