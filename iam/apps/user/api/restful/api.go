package restful

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/user"
	permission "github.com/infraboard/modules/iam/permission/restful"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

func init() {
	ioc.Api().Registry(&UserRestfulApiHandler{})
}

type UserRestfulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc user.Service
}

func (h *UserRestfulApiHandler) Name() string {
	return user.AppName
}

func (h *UserRestfulApiHandler) Init() error {
	h.svc = user.GetService()

	tags := []string{"用户管理"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET("").To(h.QueryUser).
		Doc("用户列表查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Metadata(permission.Resource("user")).
		Metadata(permission.Action("list")).
		Param(restful.QueryParameter("page_size", "页大小").DataType("integer")).
		Param(restful.QueryParameter("page_number", "页码").DataType("integer")).
		Writes(user.User{}).
		Returns(200, "OK", UserSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeUser).
		Doc("用户详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Metadata(permission.Resource("user")).
		Metadata(permission.Action("get")).
		Param(restful.PathParameter("id", "用户Id")).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}))

	ws.Route(ws.POST("").To(h.CreateUser).
		Doc("创建用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Metadata(permission.Resource("user")).
		Metadata(permission.Action("create")).
		Reads(user.CreateUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}))

	ws.Route(ws.DELETE("").To(h.DeleteUser).
		Doc("删除用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Permission(true)).
		Metadata(permission.Resource("user")).
		Metadata(permission.Action("delete")).
		Reads(user.DeleteUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}).
		Returns(404, "Not Found", nil))
	return nil
}

type UserSet struct {
	Total int64        `json:"total"`
	Items []*user.User `json:"items"`
}
