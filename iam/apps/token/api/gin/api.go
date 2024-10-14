package gin

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/modules/iam/apps/token"
)

func init() {
	ioc.Api().Registry(&TokenGinApiHandler{})
}

type TokenGinApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc token.Service
}

func (h *TokenGinApiHandler) Name() string {
	return token.AppName
}

func (h *TokenGinApiHandler) Init() error {
	h.svc = token.GetService()

	// 注册路由
	r := gin.ObjectRouter(h)
	r.POST("/", h.Login)
	r.DELETE("/", h.Logout)
	return nil
}
