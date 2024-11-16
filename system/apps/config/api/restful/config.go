package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	permission "github.com/infraboard/modules/iam/permission/restful"
	"github.com/infraboard/modules/system/apps/config"
)

func init() {
	ioc.Api().Registry(&ConfigRestfulApiHandler{})
}

type ConfigRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc config.Service
}

func (h *ConfigRestfulApiHandler) Name() string {
	return "system_config"
}

func (h *ConfigRestfulApiHandler) Init() error {
	h.svc = config.GetService()

	tags := []string{"系统配置"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryConfig).
		Doc("查询系统配置").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(ConfigSet{}).
		Returns(200, "OK", ConfigSet{}))

	ws.Route(ws.GET("/:id").To(h.DescribeConfig).
		Doc("查询配置项详情").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "ConfigItem Id")).
		Writes(config.ConfigItem{}).
		Returns(200, "OK", config.ConfigItem{}))

	ws.Route(ws.POST("").To(h.AddConfig).
		Doc("添加系统配置").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Reads(config.KVItem{}).
		Writes(config.ConfigItem{}).
		Returns(200, "OK", config.ConfigItem{}))

	ws.Route(ws.PUT("/:id").To(h.UpdateConfig).
		Doc("更新系统配置").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Param(restful.PathParameter("id", "Namespace Id")).
		Reads(config.KVItem{}).
		Writes(config.ConfigItem{}).
		Returns(200, "OK", config.ConfigItem{}))
	return nil
}

func (h *ConfigRestfulApiHandler) QueryConfig(r *restful.Request, w *restful.Response) {
	req := config.NewQueryConfigRequest()

	tk, err := h.svc.QueryConfig(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

func (h *ConfigRestfulApiHandler) DescribeConfig(r *restful.Request, w *restful.Response) {
	req := config.NewDescribeConfigRequestById(r.PathParameter("id"))

	tk, err := h.svc.DescribeConfig(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

func (h *ConfigRestfulApiHandler) AddConfig(r *restful.Request, w *restful.Response) {
	req := config.NewAddConfigRequest()
	err := r.ReadEntity(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tk, err := h.svc.AddConfig(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

func (h *ConfigRestfulApiHandler) UpdateConfig(r *restful.Request, w *restful.Response) {
	req := config.NewUpdateConfigRequest(r.PathParameter("id"))

	err := r.ReadEntity(&req.KVItem)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tk, err := h.svc.UpdateConfig(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, tk)
}

type ConfigSet struct {
	Total int64                `json:"total"`
	Items []*config.ConfigItem `json:"items"`
}
