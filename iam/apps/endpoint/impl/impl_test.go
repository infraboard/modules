package impl_test

import (
	"context"

	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/test"
)

var (
	impl endpoint.Service
	ctx  = context.Background()
)

func init() {
	test.DevelopmentSetup()
	impl = endpoint.GetService()
}
