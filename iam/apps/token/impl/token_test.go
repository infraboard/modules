package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/token"
)

func TestIssueToken(t *testing.T) {
	req := token.NewIssueTokenRequest()
	req.IssueByPassword("admin", "123456")
	req.Source = token.SOURCE_WEB
	set, err := impl.IssueToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestQueryToken(t *testing.T) {
	req := token.NewQueryTokenRequest()
	req.SetActive(true).AddUserId(1)
	set, err := impl.QueryToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestRevolkToken(t *testing.T) {
	req := token.NewRevolkTokenRequest("ZEC4VeI3dtJX9Q5Kysovaxht", "c7vJ66XYHtJ0KxhgNX4iHR8wLbGLLNtL")
	set, err := impl.RevolkToken(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}
