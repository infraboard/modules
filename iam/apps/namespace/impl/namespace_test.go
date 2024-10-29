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
