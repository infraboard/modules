package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/role"
)

func TestQueryApiPermission(t *testing.T) {
	req := role.NewQueryApiPermissionRequest()
	set, err := impl.QueryApiPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
