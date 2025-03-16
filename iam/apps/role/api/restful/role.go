package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin/binding"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/modules/iam/apps/role"
)

func (h *RoleRestfulApiHandler) QueryRole(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := role.NewQueryRoleRequest()
	if err := binding.Query.Bind(r.Request, req); err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.QueryRole(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) DescribeRole(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := role.NewDescribeRoleRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.DescribeRole(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) CreateRole(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := role.NewCreateRoleRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreateRole(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) DeleteRole(r *restful.Request, w *restful.Response) {
	req := role.NewDeleteRoleRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	u, err := h.svc.DeleteRole(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type RoleSet struct {
	Total int64        `json:"total"`
	Items []*role.Role `json:"items"`
}
