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
	req.RoleId = []uint64{1}
	set, err := impl.CreatePolicy(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestCreatePolicy2(t *testing.T) {
	req := policy.NewCreatePolicyRequest()
	// guest
	req.UserId = 3
	// 开发
	req.RoleId = []uint64{3}
	// default
	req.SetNamespaceId(1)
	// 开发小组1的资源
	req.SetScope("team", []string{"dev01"})
	set, err := impl.CreatePolicy(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
