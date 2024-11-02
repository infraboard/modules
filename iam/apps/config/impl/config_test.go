package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/config"
)

func TestQueryConfig(t *testing.T) {
	req := config.NewQueryConfigRequest()
	set, err := impl.QueryConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestDescribeConfig(t *testing.T) {
	req := config.NewDescribeConfigRequestById("1")
	ins, err := impl.DescribeConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestAddConfig(t *testing.T) {
	req := config.NewAddConfigRequest()
	req.AddKVItem()
	set, err := impl.AddConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestUpdateConfig(t *testing.T) {
	req := config.NewUpdateConfigRequest("1")
	req.Value = "xx"
	ins, err := impl.UpdateConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}
