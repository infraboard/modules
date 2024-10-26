package role

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "role"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 角色管理
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

	// 角色菜单管理
	// 添加角色关联菜单
	AddMenuPermission(context.Context, *AddMenuPermissionRequest) ([]*MenuPermission, error)
	// 移除角色关联菜单
	RemoveMenuPermission(context.Context, *RemoveMenuPermissionRequest) ([]*MenuPermission, error)
	// 更新角色权限
	UpdateMenuPermission(context.Context, *UpdateMenuPermission) ([]*MenuPermission, error)

	// 角色API接口管理
	// 添加角色关联API
	AddApiPermission(context.Context, *AddApiPermissionRequest) ([]*ApiPermission, error)
	// 移除角色关联API
	RemoveApiPermission(context.Context, *RemoveApiPermissionRequest) ([]*ApiPermission, error)
	// 更新Api权限
	UpdateApiPermission(context.Context, *UpdateApiPermissionRequest) ([]*ApiPermission, error)
}

type AddApiPermissionRequest struct {
	Items []*ApiPermissionSpec `json:"items"`
}

type RemoveApiPermissionRequest struct {
	ApiPermissionIds []uint64 `json:"api_permission_ids"`
}

type UpdateApiPermissionRequest struct {
	Items []*ApiPermission `json:"items"`
}

type AddMenuPermissionRequest struct {
	Items []*MenuPermissionSpec `json:"items"`
}

type UpdateMenuPermission struct {
	Items []MenuPermission `json:"items"`
}

type RemoveMenuPermissionRequest struct {
	MenuPermissionIds []uint64 `json:"menu_permission_ids"`
}

func NewQueryRoleRequest() *QueryRoleRequest {
	return &QueryRoleRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryRoleRequest struct {
	*request.PageRequest
	WithMenuPermission bool `json:"with_menu_permission"`
	WithApiPermission  bool `json:"with_api_permission"`
}

type DescribeRoleRequest struct {
}

type UpdateRoleRequest struct {
}

type DeleteRoleRequest struct {
}
