package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/test"
	"github.com/infraboard/modules/system/apps/config"
)

var (
	impl config.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = config.GetService()
}
