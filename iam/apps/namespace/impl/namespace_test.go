package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/namespace"
)

func TestQueryEndpoint(t *testing.T) {
	req := namespace.NewQueryNamespaceRequest()
	set, err := impl.QueryNamespace(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestCreateNamespace(t *testing.T) {
	req := namespace.NewCreateNamespaceRequest()
	req.Name = namespace.DEFAULT_NS_NAME
	req.Description = "默认空间"
	req.OwnerUserId = 3
	set, err := impl.CreateNamespace(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
