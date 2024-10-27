package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/view"
	permission "github.com/infraboard/modules/iam/permission/restful"
)

func init() {
	ioc.Api().Registry(&PolicyRestfulApiHandler{})
}

type PolicyRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc policy.PolicyService
}

func (h *PolicyRestfulApiHandler) Name() string {
	return "policy"
}

func (h *PolicyRestfulApiHandler) Init() error {
	h.svc = policy.GetService()

	tags := []string{"权限策略管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryMenu).
		Doc("策略列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(PolicySet{}).
		Returns(200, "OK", PolicySet{}))

	ws.Route(ws.GET("/:id").To(h.DescribeMenu).
		Doc("策略详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Policy Id")).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}))

	ws.Route(ws.POST("").To(h.CreateMenu).
		Doc("创建策略").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.CreateMenuRequest{}).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}))

	ws.Route(ws.DELETE("").To(h.DeleteMenu).
		Doc("删除策略").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(view.DeleteMenuRequest{}).
		Writes(view.Menu{}).
		Returns(200, "OK", view.Menu{}).
		Returns(404, "Not Found", nil))
	return nil
}

func (h *PolicyRestfulApiHandler) QueryMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewQueryPolicyRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryPolicy(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PolicyRestfulApiHandler) DescribeMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewDescribePolicyRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribePolicy(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PolicyRestfulApiHandler) CreateMenu(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := policy.NewCreatePolicyRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreatePolicy(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *PolicyRestfulApiHandler) DeleteMenu(r *restful.Request, w *restful.Response) {
	req := policy.NewDeletePolicyRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	u, err := h.svc.DeletePolicy(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type PolicySet struct {
	Total int64            `json:"total"`
	Items []*policy.Policy `json:"items"`
}
