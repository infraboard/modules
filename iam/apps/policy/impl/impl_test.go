package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl policy.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = policy.GetService()
}
