package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/policy"
)

func TestQueryNamespace(t *testing.T) {
	req := policy.NewQueryNamespaceRequest()
	req.UserId = 3
	set, err := impl.QueryNamespace(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestQueryEndpoint(t *testing.T) {
	req := policy.NewQueryEndpointRequest()
	req.UserId = 3
	req.NamespaceId = 3
	set, err := impl.QueryEndpoint(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestValidateEndpointPermission(t *testing.T) {
	req := policy.NewValidateEndpointPermissionRequest()
	req.UserId = 3
	req.SetNamespaceId(3)
	req.Service = "cmdb"
	req.Method = "GET"
	req.Path = "/cmdb/api/v1/secret"
	set, err := impl.ValidateEndpointPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestValidateEndpointPermission2(t *testing.T) {
	req := policy.NewValidateEndpointPermissionRequest()
	req.UserId = 3
	req.SetNamespaceId(1)
	req.Service = "devcloud"
	req.Method = "GET"
	req.Path = "/api/devcloud/v1/users/"
	set, err := impl.ValidateEndpointPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
