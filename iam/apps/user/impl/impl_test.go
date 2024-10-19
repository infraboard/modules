package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/user"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl user.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = user.GetService()
}
