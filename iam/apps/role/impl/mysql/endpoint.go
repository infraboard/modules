package mysql

import (
	"context"

	"github.com/infraboard/modules/iam/apps/role"
)

// 添加角色关联API
func (i *RoleServiceImpl) AddEndpiontToRole(ctx context.Context, in *role.AddEndpintToRoleRequest) ([]*role.RoleAssociateEndpointRecord, error) {
	return nil, nil
}

// 移除角色关联API
func (i *RoleServiceImpl) RemoveEndpointFromRole(ctx context.Context, in *role.RemoveEndpointFromRole) ([]*role.RoleAssociateEndpointRecord, error) {
	return nil, nil
}
