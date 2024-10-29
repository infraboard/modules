package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/role"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl role.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = role.GetService()
}
