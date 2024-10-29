package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/namespace"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl namespace.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = namespace.GetService()
}
