package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl token.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = token.GetService()
}
