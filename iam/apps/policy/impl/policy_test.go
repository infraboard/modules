package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/policy"
)

func TestQueryPolicy(t *testing.T) {
	req := policy.NewQueryPolicyRequest()
	req.WithUser = true
	req.WithRole = true
	req.WithNamespace = true
	set, err := impl.QueryPolicy(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestCreatePolicy(t *testing.T) {
	req := policy.NewCreatePolicyRequest()
	req.SetNamespaceId(3)
	req.UserId = 3
	req.RoleId = 2
	set, err := impl.CreatePolicy(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
