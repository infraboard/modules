package impl_test

import (
	"github.com/infraboard/modules/iam/test"
	"github.com/infraboard/modules/task/apps/task"
)

var (
	svc task.Service
)

func init() {
	test.DevelopmentSetup()
	svc = task.GetService()
}
