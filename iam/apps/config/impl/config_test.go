package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/config"
)

func TestQueryConfig(t *testing.T) {
	req := config.NewQueryConfigRequest()
	set, err := impl.QueryConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
