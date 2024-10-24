package policy

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "policy"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 创建策略
	CreatePolicy(context.Context, *CreatePolicyRequest) (*Policy, error)
	// 查询策略列表
	QueryPolicy(context.Context, *QueryPolicyRequest) (*types.Set[*Policy], error)
	// 查询详情
	DescribePolicy(context.Context, *DescribePolicyRequest) (*Policy, error)
	// 更新策略
	UpdatePolicy(context.Context, *UpdatePolicyRequest) (*Policy, error)
	// 删除策略
	DeletePolicy(context.Context, *DeletePolicyRequest) (*types.Set[*Policy], error)
}

func NewQueryPolicyRequest() *QueryPolicyRequest {
	return &QueryPolicyRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryPolicyRequest struct {
	*request.PageRequest
	// 关联用户Id
	UserId *uint64 `json:"user_id"`
	// 关联空间
	NamespaceId *uint64 `json:"namespace_id"`
	// 没有过期
	Expired *bool `json:"expired"`
	// 有没有启动
	Enabled *bool `json:"active"`
	// 关联查询出空间对象
	WithNamespace bool `json:"with_namespace"`
	// 关联查询出用户对象
	WithUser bool `json:"with_user"`
	// 关联查询角色对象
	WithRole bool `json:"with_role"`
}

func (r *QueryPolicyRequest) SetNamespaceId(nsId uint64) *QueryPolicyRequest {
	r.NamespaceId = &nsId
	return r
}

func (r *QueryPolicyRequest) SetUserId(uid uint64) *QueryPolicyRequest {
	r.UserId = &uid
	return r
}

func (r *QueryPolicyRequest) SetExpired(v bool) *QueryPolicyRequest {
	r.Expired = &v
	return r
}

func (r *QueryPolicyRequest) SetEnabled(v bool) *QueryPolicyRequest {
	r.Enabled = &v
	return r
}

func (r *QueryPolicyRequest) SetWithRole(v bool) *QueryPolicyRequest {
	r.WithRole = v
	return r
}
func (r *QueryPolicyRequest) SetWithUsers(v bool) *QueryPolicyRequest {
	r.WithUser = v
	return r
}
func (r *QueryPolicyRequest) SetWithUser(v bool) *QueryPolicyRequest {
	r.WithNamespace = v
	return r
}

type DescribePolicyRequest struct {
}

type UpdatePolicyRequest struct {
}

type DeletePolicyRequest struct {
}
