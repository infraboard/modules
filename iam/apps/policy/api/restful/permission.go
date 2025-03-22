package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcenter/apps/endpoint"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/view"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&PermissionRestfulApiHandler{})
}

type PermissionRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc policy.PermissionService
}

func (h *PermissionRestfulApiHandler) Name() string {
	return "permission"
}

func (h *PermissionRestfulApiHandler) Init() error {
	h.svc = policy.GetService()

	tags := []string{"用户权限查询"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("/namespaces").To(h.QueryNamespace).
		Doc("空间列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Writes(NamespaceSet{}).
		Returns(200, "OK", NamespaceSet{}))

	ws.Route(ws.GET("/menus").To(h.QueryMenu).
		Doc("菜单列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Writes(MenuSet{}).
		Returns(200, "OK", MenuSet{}))

	ws.Route(ws.GET("/endpoints").To(h.QueryEndpoint).
		Doc("接口列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Writes(EndpointSet{}).
		Returns(200, "OK", EndpointSet{}))

	return nil
}

func (h *PermissionRestfulApiHandler) QueryNamespace(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewQueryNamespaceRequest()

	// 2. 执行逻辑
	tk, err := h.svc.QueryNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PermissionRestfulApiHandler) QueryMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewQueryMenuRequest()

	// 2. 执行逻辑
	tk, err := h.svc.QueryMenu(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PermissionRestfulApiHandler) QueryEndpoint(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewQueryEndpointRequest()

	// 2. 执行逻辑
	tk, err := h.svc.QueryEndpoint(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PermissionRestfulApiHandler) ValidateEndpointPermission(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewValidateEndpointPermissionRequest()
	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.ValidateEndpointPermission(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

type MenuSet struct {
	Total int64        `json:"total"`
	Items []*view.Menu `json:"items"`
}

type NamespaceSet struct {
	Total int64                  `json:"total"`
	Items []*namespace.Namespace `json:"items"`
}

type EndpointSet struct {
	Total int64                `json:"total"`
	Items []*endpoint.Endpoint `json:"items"`
}
