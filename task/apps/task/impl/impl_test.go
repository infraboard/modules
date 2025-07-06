package impl_test

import (
	"github.com/infraboard/modules/task/apps/task"
	"github.com/infraboard/modules/task/test"
)

var (
	svc task.Service
)

func init() {
	test.DevelopmentSetup()
	svc = task.GetService()
}
