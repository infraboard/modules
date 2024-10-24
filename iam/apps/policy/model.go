package policy

import (
	"time"

	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/user"
)

func NewPolicy() *Policy {
	return &Policy{
		Meta: *apps.NewMeta(),
	}
}

type Policy struct {
	// 基础数据
	apps.Meta
	// 策略定义
	CreatePolicyRequest
	// 关联空间
	Namespace *namespace.Namespace `json:"namespace,omitempty" gorm:"-"`
	// 关联用户
	User *user.User `json:"user,omitempty" gorm:"-"`
	// 关联角色
	Role *role.Role `json:"role,omitempty" gorm:"-"`
}

func (u *Policy) TableName() string {
	return "policy"
}

func NewCreatePolicyRequest() *CreatePolicyRequest {
	return &CreatePolicyRequest{
		Extras: map[string]string{},
		Filter: map[string]string{},
	}
}

type CreatePolicyRequest struct {
	// 创建者
	CreateBy uint64 `json:"create_by" bson:"create_by" gorm:"column:create_by;type:uint" description:"创建者" optional:"true"`
	// 空间
	NamespaceId *uint64 `json:"namespace_id" bson:"namespace_id" gorm:"column:namespace_id;type:varchar(200);index" description:"策略生效的空间Id" optional:"true"`
	// 用户Id
	UserId uint64 `json:"user_id" bson:"user_id" gorm:"column:user_id;type:uint;not null;index" validate:"required" description:"被授权的用户"`
	// 角色Id
	RoleId uint64 `json:"role_id" bson:"role_id" gorm:"column:role_id;type:uint;not null;index" validate:"required" description:"被关联的角色"`
	// 过滤器
	Filter map[string]string `json:"filter" bson:"filter" gorm:"column:filter;serializer:json;type:json" description:"被关联的权限过滤器" optional:"true"`
	// 策略过期时间
	ExpiredTime *time.Time `json:"expired_time" bson:"expired_time" gorm:"column:expired_time;type:timestamp;index" description:"策略过期时间" optional:"true"`
	// 只读策略, 不允许用户修改, 一般用于系统管理
	ReadOnly bool `json:"read_only" bson:"read_only" gorm:"column:read_only;type:tinyint(1)" description:"只读策略, 不允许用户修改, 一般用于系统管理" optional:"true"`
	// 该策略是否启用
	Enabled bool `json:"enabled" bson:"enabled" gorm:"column:enabled;type:tinyint(1)" description:"该策略是否启用" optional:"true"`
	// 策略标签
	Label string `json:"label" gorm:"column:label;type:varchar(200);index" description:"策略标签" optional:"true"`
	// 扩展信息
	Extras map[string]string `json:"extras" bson:"extras" gorm:"column:extras;serializer:json;type:json" description:"扩展信息" optional:"true"`
}
