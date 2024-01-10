package impl_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/identity/apps/user"

	// 被测试对象
	_ "github.com/infraboard/modules/identity/apps/user/impl/mysql"
)

var (
	impl user.Service
	ctx  = context.Background()
)

func init() {
	ioc.DevelopmentSetup()
	impl = ioc.Controller().Get(user.AppName).(user.Service)
}
