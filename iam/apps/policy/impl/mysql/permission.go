package mysql

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/view"
)

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

// 查询用户可以访问的Api接口
func (i *PolicyServiceImpl) QueryEndpoint(ctx context.Context, in *policy.QueryEndpointRequest) (*policy.QueryEndpointResponse, error) {
	resp := policy.NewQueryEndpointResponse()
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

	// 没有权限
	if policies.Len() == 0 {
		return resp, nil
	}

	if policies.Len() > 1 {
		return nil, fmt.Errorf("同一个空间下, 一个用户有多条[%d]授权", policies.Len())
	}

	p := policies.First()
	roleReq := role.NewQueryMatchedEndpointRequest()
	roleReq.Add(p.RoleId...)

	// p.ResourceScope
	resp.ResourceScope = p.ResourceScope
	if policies.Len() > 0 {
		set, err := role.GetService().QueryMatchedEndpoint(ctx, roleReq)
		if err != nil {
			return nil, err
		}
		resp.Items = set.Items
	}

	return resp, nil
}

// 校验Api接口权限
func (i *PolicyServiceImpl) ValidateEndpointPermission(ctx context.Context, in *policy.ValidateEndpointPermissionRequest) (*policy.ValidateEndpointPermissionResponse, error) {
	resp := policy.NewValidateEndpointPermissionResponse(*in)

	// 空间Owner有所有权限
	ns, err := namespace.GetService().DescribeNamespace(ctx, namespace.NewDescribeNamespaceRequest().SetNamespaceId(in.GetNamespaceId()))
	if err != nil {
		return nil, err
	}
	if ns.IsOwner(in.UserId) {
		resp.HasPermission = true
		return resp, nil
	}

	// 非空间管理员需要独立鉴权, 查询用户可以访问的API列表
	endpointReq := policy.NewQueryEndpointRequest()
	endpointReq.UserId = in.UserId
	endpointReq.NamespaceId = in.GetNamespaceId()
	endpointSet, err := i.QueryEndpoint(ctx, endpointReq)
	if err != nil {
		return nil, err
	}
	resp.ResourceScope = endpointSet.ResourceScope

	for _, item := range endpointSet.Items {
		if item.IsMatched(in.Service, in.Method, in.Path) {
			resp.HasPermission = true
			resp.Endpoint = item
			break
		}
	}
	return resp, nil
}

// 查询用户可以访问的菜单
func (i *PolicyServiceImpl) QueryMenu(ctx context.Context, in *policy.QueryMenuRequest) (*types.Set[*view.Menu], error) {
	return nil, nil
}

// 校验Menu视图权限
func (i *PolicyServiceImpl) ValidatePagePermission(ctx context.Context, in *policy.ValidatePagePermissionRequest) (*policy.ValidatePagePermissionResponse, error) {
	return nil, nil
}
