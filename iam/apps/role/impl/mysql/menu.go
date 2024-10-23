package mysql

import (
	"context"

	"github.com/infraboard/modules/iam/apps/role"
)

// 添加角色关联菜单
func (i *RoleServiceImpl) AddMenuToRole(ctx context.Context, in *role.AddMenuToRoleRequest) ([]*role.RoleAssociateMenuRecord, error) {
	return nil, nil
}

// 移除角色关联菜单
func (i *RoleServiceImpl) RemoveMenuFromRole(ctx context.Context, in *role.RemoveMenuFromRoleRequest) ([]*role.RoleAssociateMenuRecord, error) {
	return nil, nil
}
