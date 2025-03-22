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

func TestAddViewPermission(t *testing.T) {
	req := role.NewAddViewPermissionRequest()
	req.Add(role.NewViewPermissionSpec())
	set, err := impl.AddViewPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestRemoveViewPermission(t *testing.T) {
	req := role.NewRemoveViewPermissionRequest()
	set, err := impl.RemoveViewPermission(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestQueryMatchedMenu(t *testing.T) {
	req := role.NewQueryMatchedMenuRequest()
	set, err := impl.QueryMatchedMenu(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
