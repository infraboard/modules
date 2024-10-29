package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/view"
)

func TestQueryPage(t *testing.T) {
	req := view.NewQueryPageRequest()
	set, err := impl.QueryPage(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
