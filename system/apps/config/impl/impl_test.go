package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/config"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl config.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = config.GetService()
}
