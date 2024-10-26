package mysql

import (
	"context"

	"github.com/infraboard/modules/iam/apps/role"
)

// 添加角色关联菜单
func (i *RoleServiceImpl) AddMViewPermission(ctx context.Context, in *role.AddViewPermissionRequest) ([]*role.ViewPermission, error) {
	return nil, nil
}

// 移除角色关联菜单
func (i *RoleServiceImpl) RemoveViewPermission(ctx context.Context, in *role.RemoveViewPermissionRequest) ([]*role.ViewPermission, error) {
	return nil, nil
}

// 更新角色权限
func (i *RoleServiceImpl) UpdateViewPermission(ctx context.Context, in *role.UpdateViewPermission) ([]*role.ViewPermission, error) {
	return nil, nil
}
