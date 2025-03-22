package restful

import (
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/role"
)

func (h *RoleRestfulApiHandler) QueryMatchedEndpoint(r *restful.Request, w *restful.Response) {
	v, err := strconv.ParseUint(r.PathParameter("id"), 10, 64)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}
	req := role.NewQueryMatchedEndpointRequest()
	req.Add(v)

	tk, err := h.svc.QueryMatchedEndpoint(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) AddApiPermission(r *restful.Request, w *restful.Response) {
	v, err := strconv.ParseUint(r.PathParameter("id"), 10, 64)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}
	req := role.NewAddApiPermissionRequest(v)

	err = r.ReadEntity(&req.Items)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.AddApiPermission(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, tk)
}

func (h *RoleRestfulApiHandler) RemoveApiPermission(r *restful.Request, w *restful.Response) {
	v, err := strconv.ParseUint(r.PathParameter("id"), 10, 64)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse id error, %s", err))
		return
	}
	req := role.NewRemoveApiPermissionRequest(v)

	err = r.ReadEntity(&req.ApiPermissionIds)
	if err != nil {
		response.Failed(w, err)
		return
	}

	u, err := h.svc.RemoveApiPermission(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. 返回响应
	response.Success(w, u)
}

type EndpointSet struct {
	Total int64                `json:"total"`
	Items []*endpoint.Endpoint `json:"items"`
}
