package feishu_test

import (
	"context"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/infraboard/modules/iam/apps/token/issuer/feishu"
)

var (
	client *feishu.Feishu
	ctx    = context.Background()
)

func TestGetUserInfo(t *testing.T) {
	TestLogin(t)

	u, err := client.GetUserInfo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestLogin(t *testing.T) {
	// 加载测试配置
	conf := feishu.NewDefaultConfig()
	if err := env.Parse(conf); err != nil {
		t.Fatal(err)
	}

	// 登陆
	client = feishu.NewFeishuClient(conf)
	tk, err := client.Login(ctx, "083kc3b5b13040b980a64c2947872a69")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)
}

func init() {
}