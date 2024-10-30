package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/namespace"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&NamespaceRestfulApiHandler{})
}

type NamespaceRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc namespace.Service
}

func (h *NamespaceRestfulApiHandler) Name() string {
	return "namespace"
}

func (h *NamespaceRestfulApiHandler) Init() error {
	h.svc = namespace.GetService()

	tags := []string{"空间管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryNamespace).
		Doc("空间列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(NamespaceSet{}).
		Returns(200, "OK", NamespaceSet{}))

	ws.Route(ws.GET("/:id").To(h.DescribeNamespace).
		Doc("空间详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Menu Id")).
		Writes(namespace.Namespace{}).
		Returns(200, "OK", namespace.Namespace{}))

	ws.Route(ws.POST("").To(h.CreateNamespace).
		Doc("创建空间").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(namespace.CreateNamespaceRequest{}).
		Writes(namespace.Namespace{}).
		Returns(200, "OK", namespace.Namespace{}))

	ws.Route(ws.PUT("/:id").To(h.UpdateNamespace).
		Doc("更新空间").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Namespace Id")).
		Reads(namespace.CreateNamespaceRequest{}).
		Writes(namespace.Namespace{}).
		Returns(200, "OK", namespace.Namespace{}))

	ws.Route(ws.DELETE("").To(h.DeleteNamespace).
		Doc("删除空间").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(namespace.DeleteNamespaceRequest{}).
		Writes(namespace.Namespace{}).
		Returns(200, "OK", namespace.Namespace{}).
		Returns(404, "Not Found", nil))
	return nil
}

func (h *NamespaceRestfulApiHandler) QueryNamespace(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := namespace.NewQueryNamespaceRequest()
	req.PageRequest = *request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *NamespaceRestfulApiHandler) DescribeNamespace(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := namespace.NewDescribeNamespaceRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribeNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *NamespaceRestfulApiHandler) CreateNamespace(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := namespace.NewCreateNamespaceRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreateNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *NamespaceRestfulApiHandler) UpdateNamespace(r *restful.Request, w *restful.Response) {
	req := namespace.NewUpdateNamespaceRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	err := r.ReadEntity(&req.CreateNamespaceRequest)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tk, err := h.svc.UpdateNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

func (h *NamespaceRestfulApiHandler) DeleteNamespace(r *restful.Request, w *restful.Response) {
	req := namespace.NewDeleteNamespaceRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	u, err := h.svc.DeleteNamespace(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type NamespaceSet struct {
	Total int64                  `json:"total"`
	Items []*namespace.Namespace `json:"items"`
}
