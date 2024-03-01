package gin

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
)

func init() {
	ioc.Api().Registry(&TokenApiHandler{})
}

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
