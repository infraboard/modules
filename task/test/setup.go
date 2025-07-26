package test

import (
	"github.com/infraboard/mcube/v2/ioc"

	// 被测试对象
	_ "github.com/infraboard/modules/task/apps"

	// rabbitmq总线
	_ "github.com/infraboard/mcube/v2/ioc/config/bus/rabbitmq"
)

func DevelopmentSetup() {
	ioc.DevelopmentSetup()
}
