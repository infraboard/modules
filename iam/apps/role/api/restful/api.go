package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/role"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&RoleRestfulApiHandler{})
}

type RoleRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc role.Service
}

func (h *RoleRestfulApiHandler) Name() string {
	return "role"
}

func (h *RoleRestfulApiHandler) Init() error {
	h.svc = role.GetService()

	tags := []string{"角色管理"}
	ws := gorestful.ObjectRouter(h)

	// 角色管理
	ws.Route(ws.GET("").To(h.QueryRole).
		Doc("角色列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(RoleSet{}).
		Returns(200, "OK", RoleSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeRole).
		Doc("角色详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Role Id")).
		Writes(role.Role{}).
		Returns(200, "OK", role.Role{}))

	ws.Route(ws.POST("").To(h.CreateRole).
		Doc("创建角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(role.CreateRoleRequest{}).
		Writes(role.Role{}).
		Returns(200, "OK", role.Role{}))

	ws.Route(ws.DELETE("").To(h.DeleteRole).
		Doc("删除角色").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(role.DeleteRoleRequest{}).
		Writes(role.Role{}).
		Returns(200, "OK", role.Role{}).
		Returns(404, "Not Found", nil))

	// 角色 Api权限管理
	ws.Route(ws.POST("/{id}/add_api_permission").To(h.AddApiPermission).
		Doc("添加接口访问权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Role Id")).
		Reads([]role.ApiPermissionSpec{}).
		Writes([]role.ApiPermission{}).
		Returns(200, "OK", []role.ApiPermission{}))

	ws.Route(ws.POST("/{id}/remove_api_permission").To(h.RemoveApiPermission).
		Doc("移除接口访问权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads([]uint64{}).
		Writes([]role.ApiPermission{}).
		Returns(200, "OK", []role.ApiPermission{}))

	ws.Route(ws.GET("/{id}/endpoints").To(h.QueryMatchedEndpoint).
		Doc("查询允许访问的接口").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Writes(EndpointSet{}).
		Returns(200, "OK", EndpointSet{}))

	// 角色视图权限管理
	ws.Route(ws.POST("/{id}/add_view_permission").To(h.AddViewPermission).
		Doc("添加视图访问权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Role Id")).
		Reads([]*role.ViewPermissionSpec{}).
		Writes([]role.ViewPermission{}).
		Returns(200, "OK", []role.ViewPermission{}))

	ws.Route(ws.POST("/{id}/remove_view_permission").To(h.RemoveViewPermission).
		Doc("移除视图访问权限").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads([]uint64{}).
		Writes([]*role.ViewPermission{}).
		Returns(200, "OK", []role.ViewPermission{}))

	ws.Route(ws.GET("/{id}/menus").To(h.QueryMatchedMenu).
		Doc("查询允许访问的菜单").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Writes(MenuSet{}).
		Returns(200, "OK", MenuSet{}))

	return nil
}
