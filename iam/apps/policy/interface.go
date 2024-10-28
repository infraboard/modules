package policy

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/view"
)

const (
	AppName = "policy"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 策略管理
	PolicyService
	// 权限查询, 整合用户多个角色的权限合集
	PermissionService
}

type PolicyService interface {
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

func NewDescribePolicyRequest() *DescribePolicyRequest {
	return &DescribePolicyRequest{}
}

type DescribePolicyRequest struct {
	apps.GetRequest
}

type UpdatePolicyRequest struct {
}

func NewDeletePolicyRequest() *DeletePolicyRequest {
	return &DeletePolicyRequest{}
}

type DeletePolicyRequest struct {
	apps.GetRequest
}

type PermissionService interface {
	// 查询用户可以访问的空间
	QueryNamespace(context.Context, *QueryNamespaceRequest) (*types.Set[*namespace.Namespace], error)
	// 查询用户可以访问的菜单
	QueryMenu(context.Context, *QueryMenuRequest) (*types.Set[*view.Menu], error)
	// 查询用户可以访问的Api接口
	QueryEndpoint(context.Context, *QueryEndpointRequest) (*types.Set[*endpoint.Endpoint], error)
}

func NewQueryNamespaceRequest() *QueryNamespaceRequest {
	return &QueryNamespaceRequest{}
}

type QueryNamespaceRequest struct {
	UserId      uint64 `json:"user_id"`
	NamespaceId uint64 `json:"namespace_id"`
}

func (r *QueryNamespaceRequest) SetUserId(v uint64) *QueryNamespaceRequest {
	r.UserId = v
	return r
}

func (r *QueryNamespaceRequest) SetNamespaceId(v uint64) *QueryNamespaceRequest {
	r.NamespaceId = v
	return r
}

func NewQueryMenuRequest() *QueryMenuRequest {
	return &QueryMenuRequest{}
}

type QueryMenuRequest struct {
}

func NewQueryEndpointRequest() *QueryEndpointRequest {
	return &QueryEndpointRequest{}
}

type QueryEndpointRequest struct {
}
