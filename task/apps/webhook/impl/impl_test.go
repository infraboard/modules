package impl_test

import (
	"github.com/infraboard/modules/task/apps/webhook"
	"github.com/infraboard/modules/task/test"
)

var (
	svc webhook.Service
)

func init() {
	test.DevelopmentSetup()
	svc = webhook.GetService()
}
