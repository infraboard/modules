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
	AddMenuToRole(context.Context, *AddMenuToRoleRequest) ([]*RoleAssociateMenuRecord, error)
	// 移除角色关联菜单
	RemoveMenuFromRole(context.Context, *RemoveMenuFromRoleRequest) ([]*RoleAssociateMenuRecord, error)

	// 角色API接口管理
	// 添加角色关联API
	AddEndpiontToRole(context.Context, *AddEndpintToRoleRequest) ([]*RoleAssociateEndpointRecord, error)
	// 移除角色关联API
	RemoveEndpointFromRole(context.Context, *RemoveEndpointFromRole) ([]*RoleAssociateEndpointRecord, error)
}

type AddEndpintToRoleRequest struct {
	Items []*RoleAssociateEndpointRecord `json:"items"`
}

type RemoveEndpointFromRole struct {
	RecordIds []uint64 `json:"record_ids"`
}

type AddMenuToRoleRequest struct {
	Items []*RoleAssociateMenuRecord `json:"items"`
}

type RemoveMenuFromRoleRequest struct {
	RecordIds []uint64 `json:"record_ids"`
}

func NewQueryRoleRequest() *QueryRoleRequest {
	return &QueryRoleRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryRoleRequest struct {
	*request.PageRequest
	WithMenu     bool `json:"with_menu"`
	WithEndpoint bool `json:"with_endpoint"`
}

type DescribeRoleRequest struct {
}

type UpdateRoleRequest struct {
}

type DeleteRoleRequest struct {
}
