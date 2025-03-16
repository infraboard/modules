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

func TestRegistryEndpoint(t *testing.T) {
	req := endpoint.NewRegistryEndpointRequest()
	re := endpoint.NewRouteEntry()
	re.Service = "cmdb"
	re.Method = "GET"
	re.Path = "/cmdb/api/v1/secret"
	re.Resource = "secret"
	re.Action = "list"
	req.AddItem(re)
	set, err := impl.RegistryEndpoint(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
