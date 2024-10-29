package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/endpoint"
)

func TestQueryEndpoint(t *testing.T) {
	req := endpoint.NewQueryEndpointRequest()
	set, err := impl.QueryEndpoint(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
