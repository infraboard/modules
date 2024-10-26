package mysql

import (
	"context"

	"github.com/infraboard/modules/iam/apps/role"
)

// 添加角色关联菜单
func (i *RoleServiceImpl) AddMenuPermission(ctx context.Context, in *role.AddMenuPermissionRequest) ([]*role.MenuPermission, error) {
	return nil, nil
}

// 移除角色关联菜单
func (i *RoleServiceImpl) RemoveMenuPermission(ctx context.Context, in *role.RemoveMenuPermissionRequest) ([]*role.MenuPermission, error) {
	return nil, nil
}

// 更新角色权限
func (i *RoleServiceImpl) UpdateMenuPermission(ctx context.Context, in *role.UpdateMenuPermission) ([]*role.MenuPermission, error) {
	return nil, nil
}
