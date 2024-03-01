package impl_test

import (
	"context"
	"os"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/user"

	// 被测试对象
	_ "github.com/infraboard/modules/iam/apps/user/impl/mysql"
)

var (
	impl user.Service
	ctx  = context.Background()
)

func init() {
	os.Setenv("DATASOURCE_DB", "test")
	os.Setenv("DATASOURCE_USERNAME", "root")
	os.Setenv("DATASOURCE_PASSWORD", "123456")

	ioc.DevelopmentSetup()
	impl = ioc.Controller().Get(user.AppName).(user.Service)
}
