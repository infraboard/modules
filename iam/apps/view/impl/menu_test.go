package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/view"
)

func TestQueryMenu(t *testing.T) {
	req := view.NewQueryMenuRequest()
	set, err := impl.QueryMenu(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestCreateMenu(t *testing.T) {
	req := view.NewCreateMenuRequest()
	req.Path = "/system/develop/tool"
	req.Name = "研发工具"
	set, err := impl.CreateMenu(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
