package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/view"
	"gorm.io/gorm"
)

// 添加角色关联菜单
func (i *RoleServiceImpl) AddViewPermission(ctx context.Context, in *role.AddViewPermissionRequest) ([]*role.ViewPermission, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate add view permission error, %s", err)
	}

	perms := []*role.ViewPermission{}
	if err := datasource.DBFromCtx(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range in.Items {
			item := in.Items[i]
			perm := role.NewViewPermission(in.RoleId, item)
			if err := tx.Save(perm).Error; err != nil {
				return err
			}
			perms = append(perms, perm)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return perms, nil
}

// 查询角色关联的视图权限
func (i *RoleServiceImpl) QueryViewPermission(ctx context.Context, in *role.QueryViewPermissionRequest) ([]*role.ViewPermission, error) {
	perms := []*role.ViewPermission{}
	if err := datasource.DBFromCtx(ctx).
		Model(&role.ViewPermission{}).
		Order("created_at desc").
		Where("id IN ?", in.RoleIds).
		Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
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
