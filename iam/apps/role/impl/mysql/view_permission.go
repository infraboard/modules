package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/view"
)

// 添加角色关联菜单
func (i *RoleServiceImpl) AddViewPermission(ctx context.Context, in *role.AddViewPermissionRequest) ([]*role.ViewPermission, error) {
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

// 查询能匹配到视图菜单
func (i *RoleServiceImpl) QueryMatchedMenu(ctx context.Context, in *role.QueryMatchedMenuRequest) (*types.Set[*view.Menu], error) {
	return nil, nil
}
