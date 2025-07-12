package impl_test

import (
	"github.com/infraboard/modules/task/apps/cronjob"
	"github.com/infraboard/modules/task/test"
)

var (
	svc cronjob.Service
)

func init() {
	test.DevelopmentSetup()
	svc = cronjob.GetService()
}
