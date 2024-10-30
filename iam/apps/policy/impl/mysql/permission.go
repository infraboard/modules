package mysql

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/view"
)

// 查询用户可以访问的Api接口
func (i *PolicyServiceImpl) QueryEndpoint(ctx context.Context, in *policy.QueryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	policies, err := i.QueryPolicy(ctx,
		policy.NewQueryPolicyRequest().
			SetSkipPage(true).
			SetNamespaceId(in.NamespaceId).
			SetUserId(in.UserId).
			SetExpired(false).
			SetEnabled(true))
	if err != nil {
		return nil, err
	}

	fmt.Println(policies)
	return nil, nil
}

// 查询用户可以访问的菜单
func (i *PolicyServiceImpl) QueryMenu(ctx context.Context, in *policy.QueryMenuRequest) (*types.Set[*view.Menu], error) {
	return nil, nil
}

// 查询用户可以访问的空间
func (i *PolicyServiceImpl) QueryNamespace(ctx context.Context, in *policy.QueryNamespaceRequest) (*types.Set[*namespace.Namespace], error) {
	nsReq := namespace.NewQueryNamespaceRequest()

	policies, err := i.QueryPolicy(ctx,
		policy.NewQueryPolicyRequest().
			SetSkipPage(true).
			SetUserId(in.UserId).
			SetExpired(false).
			SetEnabled(true))
	if err != nil {
		return nil, err
	}

	policies.ForEach(func(t *policy.Policy) {
		if t.NamespaceId != nil {
			nsReq.AddNamespaceIds(*t.NamespaceId)
		}
	})

	return i.namespace.QueryNamespace(ctx, nsReq)
}
