package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/role"
)

// 添加角色关联API
func (i *RoleServiceImpl) AddApiPermission(ctx context.Context, in *role.AddApiPermissionRequest) ([]*role.ApiPermission, error) {
	return nil, nil
}

// 移除角色关联API
func (i *RoleServiceImpl) RemoveApiPermission(ctx context.Context, in *role.RemoveApiPermissionRequest) ([]*role.ApiPermission, error) {
	return nil, nil
}

// 更新Api权限
func (i *RoleServiceImpl) UpdateApiPermission(ctx context.Context, in *role.UpdateApiPermissionRequest) ([]*role.ApiPermission, error) {
	return nil, nil
}

// 查询匹配到的Api接口列表
func (i *RoleServiceImpl) QueryMatchedEndpoint(ctx context.Context, in *role.QueryMatchedEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	return nil, nil
}
