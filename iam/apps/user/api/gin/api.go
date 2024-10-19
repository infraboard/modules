package gin

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/apps/user"
	permission "github.com/infraboard/modules/iam/permission/gin"
)

func init() {
	ioc.Api().Registry(&UserGinApiHandler{})
}

type UserGinApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc user.Service
}

func (h *UserGinApiHandler) Name() string {
	return user.AppName
}

func (h *UserGinApiHandler) Init() error {
	h.svc = user.GetService()

	// 注册路由
	r := gin.ObjectRouter(h)
	r.Use(
		// 认证
		permission.Auth(),
		// 鉴权
		permission.Required(role.ADMIN),
	)

	r.GET("/", h.QueryUser)
	r.GET("/:id", h.DescribeUser)
	r.POST("/", h.CreateUser)
	r.DELETE("/:id", h.DeleteUser)
	return nil
}
