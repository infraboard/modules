package test

import (
	"os"

	"github.com/infraboard/mcube/v2/ioc"

	// 被测试对象
	_ "github.com/infraboard/modules/iam/init/mysql"
)

func DevelopmentSetup() {
	os.Setenv("DATASOURCE_DB", "test")
	os.Setenv("DATASOURCE_USERNAME", "root")
	os.Setenv("DATASOURCE_PASSWORD", "123456")
	// os.Setenv("DATASOURCE_AUTO_MIGRATE", "true")
	ioc.DevelopmentSetup()
}
