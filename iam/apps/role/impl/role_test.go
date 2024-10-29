package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/role"
)

func TestQueryRole(t *testing.T) {
	req := role.NewQueryRoleRequest()
	set, err := impl.QueryRole(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
