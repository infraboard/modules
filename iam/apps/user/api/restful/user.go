package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/modules/iam/apps/user"
)

func (h *UserRestfulApiHandler) QueryUser(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewQueryUserRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(r.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *UserRestfulApiHandler) DescribeUser(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewDescribeUserRequestById(r.PathParameter("id"))

	// 2. 执行逻辑
	tk, err := h.svc.DescribeUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *UserRestfulApiHandler) CreateUser(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewCreateUserRequest()

	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreateUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *UserRestfulApiHandler) DeleteUser(r *restful.Request, w *restful.Response) {
	req := user.NewDeleteUserRequest(r.PathParameter("id"))
	u, err := h.svc.DeleteUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}
