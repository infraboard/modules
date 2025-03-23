package restful

import (
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/view"
)

func (h *RoleRestfulApiHandler) QueryMatchedPage(r *restful.Request, w *restful.Response) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := role.NewQueryMatchedPageRequest()
	if err := req.SetIdByString(r.PathParameter("id")); err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.QueryMatchedPage(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) AddViewPermission(r *restful.Request, w *restful.Response) {
	req := role.NewAddViewPermissionRequest()
	v, err := strconv.ParseUint(r.PathParameter("id"), 10, 64)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}
	req.RoleId = v

	err = r.ReadEntity(&req.Items)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.AddViewPermission(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) RemoveViewPermission(r *restful.Request, w *restful.Response) {
	req := role.NewRemoveViewPermissionRequest()
	v, err := strconv.ParseUint(r.PathParameter("id"), 10, 64)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}
	req.RoleId = v

	err = r.ReadEntity(&req.ViewPermissionIds)
	if err != nil {
		response.Failed(w, err)
		return
	}

	u, err := h.svc.RemoveViewPermission(r.Request.Context(), req)
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
