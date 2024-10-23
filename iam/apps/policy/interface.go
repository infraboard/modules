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
	CreatePolicy(context.Context, CreatePolicyRequest) (*Policy, error)
	// 查询策略列表
	QueryPolicy(context.Context, QueryPolicyRequest) (*types.Set[*Policy], error)
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
}

type DescribePolicyRequest struct {
}

type UpdatePolicyRequest struct {
}

type DeletePolicyRequest struct {
}
