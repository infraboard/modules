package role

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "roles"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 创建角色
	CreateRole(context.Context, *CreateRoleRequest) (*Role, error)
	// 列表查询
	QueryRole(context.Context, *QueryRoleRequest) (*types.Set[*Role], error)
	// 详情查询
	DescribeRole(context.Context, *DescribeRoleRequest) (*Role, error)
	// 更新角色
	UpdateRole(context.Context, *UpdateRoleRequest) (*Role, error)
	// 删除角色
	DeleteRole(context.Context, *DeleteRoleRequest) (*Role, error)
}

type QueryRoleRequest struct {
}

type DescribeRoleRequest struct {
}

type UpdateRoleRequest struct {
}

type DeleteRoleRequest struct {
}
