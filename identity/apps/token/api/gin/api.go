package gin

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/identity/apps/token"
)

func init() {
	ioc.Api().Registry(&TokenApiHandler{})
}

// 不适用接口, 直接定义Gin的一个handlers
// 什么是Gin的Handler  HandlerFunc
// HandlerFunc defines the handler used by gin middleware as return value.
// type HandlerFunc func(*Context)
// HandleFunc 只是定义 如何处理 HTTP 的请求与响应

type TokenApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc token.Service
}

func (t *TokenApiHandler) Name() string {
	return token.AppName
}

func (t *TokenApiHandler) Init() error {
	t.svc = ioc.Controller().Get(token.AppName).(token.Service)
	return nil
}
