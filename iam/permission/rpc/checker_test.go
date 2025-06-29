package permission_test

import (
	"context"
	"testing"

	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/token"
	permission "github.com/infraboard/modules/iam/permission/rpc"
	"github.com/infraboard/modules/iam/test"
)

func TestValiateToken(t *testing.T) {
	tk, err := permission.GetPermissionChecker().
		ValiateToken(context.Background(), token.NewValiateTokenRequest("fY2dsHqkSBNATR1jTHsXIpXo"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
}

func TestValidateEndpointPermission(t *testing.T) {
	req := policy.NewValidateEndpointPermissionRequest()
	req.UserId = 5
	req.SetNamespaceId(1)
	req.Service = "cmdb"
	req.Method = "GET"
	req.Path = "/api/cmdb/v1/secret"
	resp, err := permission.GetPermissionChecker().
		ValidateEndpointPermission(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestRegistryEndpoint(t *testing.T) {
	req := endpoint.NewRegistryEndpointRequest()
	req.AddItem(&endpoint.RouteEntry{
		Service:  "cmdb",
		Resource: "secret",
		Action:   "delete",
		Method:   "DELETE",
		Path:     "/api/cmdb/v1/secret/{id}",
	})
	resp, err := permission.GetApiRegister().
		RegistryEndpoint(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func init() {
	test.DevelopmentSetup()
}
