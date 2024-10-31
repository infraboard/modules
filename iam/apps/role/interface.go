package role

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/view"
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

func NewDescribeRoleRequest() *DescribeRoleRequest {
	return &DescribeRoleRequest{}
}

type DescribeRoleRequest struct {
	apps.GetRequest
}

type UpdateRoleRequest struct {
	apps.GetRequest
	CreateRoleRequest
}

func NewDeleteRoleRequest() *DeleteRoleRequest {
	return &DeleteRoleRequest{}
}

type DeleteRoleRequest struct {
	apps.GetRequest
}

// 角色API接口管理
type ApiPermissionService interface {
	// 查询角色关联的权限条目
	QueryApiPermission(context.Context, *QueryApiPermissionRequest) ([]*ApiPermission, error)
	// 添加角色关联API
	AddApiPermission(context.Context, *AddApiPermissionRequest) ([]*ApiPermission, error)
	// 移除角色关联API
	RemoveApiPermission(context.Context, *RemoveApiPermissionRequest) ([]*ApiPermission, error)
	// 查询匹配到的Api接口列表
	QueryMatchedEndpoint(context.Context, *QueryMatchedEndpointRequest) (*types.Set[*endpoint.Endpoint], error)
}

func NewQueryApiPermissionRequest() *QueryApiPermissionRequest {
	return &QueryApiPermissionRequest{
		RoleIds:          []uint64{},
		ApiPermissionIds: []uint64{},
	}
}

type QueryApiPermissionRequest struct {
	RoleIds          []uint64 `json:"role_ids"`
	ApiPermissionIds []uint64 `json:"api_permission_ids"`
}

func (r *QueryApiPermissionRequest) AddRoleId(roleIds ...uint64) *QueryApiPermissionRequest {
	r.RoleIds = append(r.RoleIds, roleIds...)
	return r
}

func (r *QueryApiPermissionRequest) AddPermissionId(permissionIds ...uint64) *QueryApiPermissionRequest {
	r.ApiPermissionIds = append(r.ApiPermissionIds, permissionIds...)
	return r
}

func NewQueryMatchedEndpointRequest() *QueryMatchedEndpointRequest {
	return &QueryMatchedEndpointRequest{}
}

type QueryMatchedEndpointRequest struct {
	apps.GetRequest
}

func NewAddApiPermissionRequest() *AddApiPermissionRequest {
	return &AddApiPermissionRequest{}
}

type AddApiPermissionRequest struct {
	RoleId uint64               `json:"role_id"`
	Items  []*ApiPermissionSpec `json:"items"`
}

func (r *AddApiPermissionRequest) Validate() error {
	return validator.Validate(r)
}

func NewRemoveApiPermissionRequest() *RemoveApiPermissionRequest {
	return &RemoveApiPermissionRequest{
		ApiPermissionIds: []uint64{},
	}
}

type RemoveApiPermissionRequest struct {
	RoleId           uint64   `json:"role_id"`
	ApiPermissionIds []uint64 `json:"api_permission_ids"`
}

func (r *RemoveApiPermissionRequest) Validate() error {
	return validator.Validate(r)
}

type UpdateApiPermissionRequest struct {
	Items []*ApiPermission `json:"items"`
}

// 角色菜单管理
type ViewPermissionService interface {
	// 查询角色关联的视图权限
	QueryViewPermission(context.Context, *QueryViewPermissionRequest) ([]*ViewPermission, error)
	// 添加角色关联菜单
	AddViewPermission(context.Context, *AddViewPermissionRequest) ([]*ViewPermission, error)
	// 移除角色关联菜单
	RemoveViewPermission(context.Context, *RemoveViewPermissionRequest) ([]*ViewPermission, error)
	// 查询能匹配到视图菜单
	QueryMatchedMenu(context.Context, *QueryMatchedMenuRequest) (*types.Set[*view.Menu], error)
}

func NewQueryViewPermissionRequest() *QueryViewPermissionRequest {
	return &QueryViewPermissionRequest{
		RoleIds:           []uint64{},
		ViewPermissionIds: []uint64{},
	}
}

type QueryViewPermissionRequest struct {
	RoleIds           []uint64 `json:"role_ids"`
	ViewPermissionIds []uint64 `json:"view_permission_ids"`
}

func (r *QueryViewPermissionRequest) AddRoleId(roleIds ...uint64) *QueryViewPermissionRequest {
	r.RoleIds = append(r.RoleIds, roleIds...)
	return r
}

func (r *QueryViewPermissionRequest) AddPermissionId(permissionIds ...uint64) *QueryViewPermissionRequest {
	r.ViewPermissionIds = append(r.ViewPermissionIds, permissionIds...)
	return r
}

func NewQueryMatchedMenuRequest() *QueryMatchedMenuRequest {
	return &QueryMatchedMenuRequest{}
}

type QueryMatchedMenuRequest struct {
	apps.GetRequest
}

func NewAddViewPermissionRequest() *AddViewPermissionRequest {
	return &AddViewPermissionRequest{
		Items: []*ViewPermissionSpec{},
	}
}

type AddViewPermissionRequest struct {
	RoleId uint64                `json:"role_id"`
	Items  []*ViewPermissionSpec `json:"items"`
}

func (r *AddViewPermissionRequest) Validate() error {
	return validator.Validate(r)
}

type UpdateViewPermission struct {
	Items []ViewPermission `json:"items"`
}

func NewRemoveViewPermissionRequest() *RemoveViewPermissionRequest {
	return &RemoveViewPermissionRequest{
		ViewPermissionIds: []uint64{},
	}
}

type RemoveViewPermissionRequest struct {
	RoleId            uint64   `json:"role_id"`
	ViewPermissionIds []uint64 `json:"menu_permission_ids"`
}

func (r *RemoveViewPermissionRequest) Validate() error {
	return validator.Validate(r)
}
