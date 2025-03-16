package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/view"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&PageRestfulApiHandler{})
}

type PageRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc view.PageService
}

func (h *PageRestfulApiHandler) Name() string {
	return "page"
}

func (h *PageRestfulApiHandler) Init() error {
	h.svc = view.GetService()

	tags := []string{"页面管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryPage).
		Doc("页面列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(PageSet{}).
		Returns(200, "OK", PageSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribePage).
		Doc("页面详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Page Id")).
		Writes(view.Page{}).
		Returns(200, "OK", view.Page{}))

	ws.Route(ws.POST("").To(h.CreatePage).
		Doc("创建页面").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.CreateMenuRequest{}).
		Writes(view.Page{}).
		Returns(200, "OK", view.Page{}))

	ws.Route(ws.DELETE("").To(h.DeletePage).
		Doc("删除页面").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.DeleteMenuRequest{}).
		Writes(view.Page{}).
		Returns(200, "OK", view.Page{}).
		Returns(404, "Not Found", nil))
	return nil
}

func (h *PageRestfulApiHandler) QueryPage(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewQueryPageRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryPage(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PageRestfulApiHandler) DescribePage(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewDescribePageRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribePage(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PageRestfulApiHandler) CreatePage(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewCreatePageRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreatePage(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PageRestfulApiHandler) DeletePage(r *restful.Request, w *restful.Response) {
	req := view.NewDeletePageRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	u, err := h.svc.DeletePage(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type PageSet struct {
	Total int64        `json:"total"`
	Items []*view.Page `json:"items"`
}
