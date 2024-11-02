package impl_test

import (
	"testing"

	"github.com/infraboard/modules/iam/apps/config"
	"github.com/infraboard/modules/iam/apps/token/issuer/ldap"
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
	req := config.NewDescribeConfigRequestByKey("ldap")
	ins, err := impl.DescribeConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	ldapConf := ldap.NewConfig()
	if err := ins.Load(ldapConf); err != nil {
		t.Fatal(err)
	}
	t.Log(ldapConf)
}

func TestAddConfig(t *testing.T) {
	req := config.NewAddConfigRequest()
	ldapConf := ldap.NewConfig()
	ldapConf.BindDn = "cn=admin,dc=example,dc=org"
	ldapConf.BindPassword = "123456"
	req.AddKVItem(config.NewKVItem("ldap", ldapConf.String()).SetGroup("登录设置").SetIsEncrypted(true))
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
