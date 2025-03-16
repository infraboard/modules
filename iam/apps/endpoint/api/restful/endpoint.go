package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/endpoint"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&EndpointRestfulApiHandler{})
}

type EndpointRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc endpoint.Service
}

func (h *EndpointRestfulApiHandler) Name() string {
	return "endpoint"
}

func (h *EndpointRestfulApiHandler) Init() error {
	h.svc = endpoint.GetService()

	tags := []string{"API管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryEndpoint).
		Doc("API列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(EndpointSet{}).
		Returns(200, "OK", EndpointSet{}))

	ws.Route(ws.GET("/:id").To(h.DescribeEndpoint).
		Doc("API详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Menu Id")).
		Writes(endpoint.Endpoint{}).
		Returns(200, "OK", endpoint.Endpoint{}))

	ws.Route(ws.POST("").To(h.RegistryEndpoint).
		Doc("注册API").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads([]endpoint.RouteEntry{}).
		Writes(EndpointSet{}).
		Returns(200, "OK", EndpointSet{}))
	return nil
}

func (h *EndpointRestfulApiHandler) QueryEndpoint(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := endpoint.NewQueryEndpointRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryEndpoint(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *EndpointRestfulApiHandler) DescribeEndpoint(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := endpoint.NewDescribeEndpointRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribeEndpoint(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *EndpointRestfulApiHandler) RegistryEndpoint(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := endpoint.NewRegistryEndpointRequest()

	err := r.ReadEntity(&req.Items)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.RegistryEndpoint(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

type EndpointSet struct {
	Total int64                `json:"total"`
	Items []*endpoint.Endpoint `json:"items"`
}
