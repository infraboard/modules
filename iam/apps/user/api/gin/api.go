package gin

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/user"
)

func init() {
	ioc.Api().Registry(&UserApiHandler{})
}

type UserApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc user.Service
}

func (t *UserApiHandler) Name() string {
	return user.AppName
}

func (t *UserApiHandler) Init() error {
	t.svc = ioc.Controller().Get(user.AppName).(user.Service)
	return nil
}
