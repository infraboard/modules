package mongo_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/oss"
)

var (
	impl oss.Service
	ctx  = context.Background()
)

func init() {
	ioc.DevelopmentSetup()
	impl = oss.GetService()
}
