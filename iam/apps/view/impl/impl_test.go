package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/view"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl view.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = view.GetService()
}
