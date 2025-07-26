package impl_test

import (
	"github.com/infraboard/modules/task/apps/event"
	"github.com/infraboard/modules/task/test"
)

var (
	svc event.Service
)

func init() {
	test.DevelopmentSetup()
	svc = event.GetService()
}
