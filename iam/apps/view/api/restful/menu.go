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
	ioc.Api().Registry(&MenuRestfulApiHandler{})
}

type MenuRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc view.MenuService
}

func (h *MenuRestfulApiHandler) Name() string {
	return "menu"
}

func (h *MenuRestfulApiHandler) Init() error {
	h.svc = view.GetService()

	tags := []string{"菜单管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryMenu).
		Doc("菜单列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(MenuSet{}).
		Returns(200, "OK", MenuSet{}))

	ws.Route(ws.GET("/:id").To(h.DescribeMenu).
		Doc("菜单详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Menu Id")).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}))

	ws.Route(ws.POST("").To(h.CreateMenu).
		Doc("创建菜单").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.CreateMenuRequest{}).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}))

	ws.Route(ws.DELETE("").To(h.DeleteMenu).
		Doc("删除菜单").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.DeleteMenuRequest{}).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}).
		Returns(404, "Not Found", nil))
	return nil
}

func (h *MenuRestfulApiHandler) QueryMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewQueryMenuRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryMenu(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *MenuRestfulApiHandler) DescribeMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewDescribeMenuRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribeMenu(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *MenuRestfulApiHandler) CreateMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := view.NewCreateMenuRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreateMenu(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *MenuRestfulApiHandler) DeleteMenu(r *restful.Request, w *restful.Response) {
	req := view.NewDeleteMenuRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	u, err := h.svc.DeleteMenu(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type MenuSet struct {
	Total int64        `json:"total"`
	Items []*view.Menu `json:"items"`
}
