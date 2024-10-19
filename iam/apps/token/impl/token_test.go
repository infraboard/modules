package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/token"
)

func TestIssueToken(t *testing.T) {
	req := token.NewIssueTokenRequest()
	req.IssueByPassword("admin", "123456")
	set, err := impl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestQueryToken(t *testing.T) {
	req := token.NewQueryTokenRequest()
	set, err := impl.QueryToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
