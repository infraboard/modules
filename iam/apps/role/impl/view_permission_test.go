package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/role"
)

func TestQueryViewPermission(t *testing.T) {
	req := role.NewQueryViewPermissionRequest()
	set, err := impl.QueryViewPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
