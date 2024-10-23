package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/role"
)

// 创建角色
func (i *RoleServiceImpl) CreateRole(ctx context.Context, in *role.CreateRoleRequest) (*role.Role, error) {
	return nil, nil
}

// 列表查询
func (i *RoleServiceImpl) QueryRole(ctx context.Context, in *role.QueryRoleRequest) (*types.Set[*role.Role], error) {
	return nil, nil
}

// 详情查询
func (i *RoleServiceImpl) DescribeRole(ctx context.Context, in *role.DescribeRoleRequest) (*role.Role, error) {
	return nil, nil
}

// 更新角色
func (i *RoleServiceImpl) UpdateRole(ctx context.Context, in *role.UpdateRoleRequest) (*role.Role, error) {
	return nil, nil
}

// 删除角色
func (i *RoleServiceImpl) DeleteRole(ctx context.Context, in *role.DeleteRoleRequest) (*role.Role, error) {
	return nil, nil
}
