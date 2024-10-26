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
	RoleService
	ApiPermissionService
	ViewPermissionService
}

// 角色管理
type RoleService interface {
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

// 角色API接口管理
type ApiPermissionService interface {
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

// 角色菜单管理
type ViewPermissionService interface {
	// 添加角色关联菜单
	AddMViewPermission(context.Context, *AddViewPermissionRequest) ([]*ViewPermission, error)
	// 移除角色关联菜单
	RemoveViewPermission(context.Context, *RemoveViewPermissionRequest) ([]*ViewPermission, error)
	// 更新角色权限
	UpdateViewPermission(context.Context, *UpdateViewPermission) ([]*ViewPermission, error)
}

type AddViewPermissionRequest struct {
	Items []*ViewPermissionSpec `json:"items"`
}

type UpdateViewPermission struct {
	Items []ViewPermission `json:"items"`
}

type RemoveViewPermissionRequest struct {
	ViewPermissionIds []uint64 `json:"menu_permission_ids"`
}
