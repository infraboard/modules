package restful

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/user"
	permission "github.com/infraboard/modules/iam/permission/restful"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
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
		Reads(user.QueryUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", types.New[*user.User]()))

	ws.Route(ws.GET("/:id").To(h.DescribeUser).
		Doc("用户详情查询").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Reads(user.DescribeUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}))

	ws.Route(ws.POST("").To(h.CreateUser).
		Doc("创建用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Reads(user.CreateUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}))

	ws.Route(ws.DELETE("").To(h.DeleteUser).
		Doc("删除用户").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(permission.Auth(true)).
		Metadata(permission.Required(user.ROLE_ADMIN)).
		Reads(user.DeleteUserRequest{}).
		Writes(user.User{}).
		Returns(200, "OK", user.User{}))

	// ws.Route(ws.DELETE("/").To(h.Logout).
	// 	Doc("验证令牌").
	// 	Metadata(restfulspec.KeyOpenAPITags, tags).
	// 	Metadata(label.Auth, label.Enable).
	// 	Metadata(label.PERMISSION_MODE, label.PERMISSION_MODE_ACL.Value()).
	// 	Reads(token.LoginRequest{}).
	// 	Writes(token.Token{}).
	// 	Returns(200, "OK", token.Token{}).
	// 	Returns(404, "Not Found", nil))

	return nil
}
