package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/role"
	"gorm.io/gorm"
)

// 添加角色关联API
func (i *RoleServiceImpl) AddApiPermission(ctx context.Context, in *role.AddApiPermissionRequest) ([]*role.ApiPermission, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate add api permission error, %s", err)
	}

	perms := []*role.ApiPermission{}
	if err := datasource.DBFromCtx(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range in.Items {
			item := in.Items[i]
			perm := role.NewApiPermission(in.RoleId, item)
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

// 移除角色关联API
func (i *RoleServiceImpl) QueryApiPermission(ctx context.Context, in *role.QueryApiPermissionRequest) ([]*role.ApiPermission, error) {
	perms := []*role.ApiPermission{}
	if err := datasource.DBFromCtx(ctx).
		Model(&role.ApiPermission{}).
		Order("created_at desc").
		Where("id IN ?", in.RoleIds).
		Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
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
