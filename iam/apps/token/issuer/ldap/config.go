package ldap

import "github.com/infraboard/mcube/v2/tools/pretty"

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	// 开启LDAP认证
	Enabled bool `json:"enabled" toml:"enabled" yaml:"enabled" env:"ENABLED"`
	// LDAP Server URL
	Url string `json:"url" toml:"url" yaml:"url" env:"URL"`
	// 管理账号的用户名称
	BindDn string `json:"bind_dn" toml:"bind_dn" yaml:"bind_dn" env:"BIND_DN"`
	// 管理账号的用户密码
	BindPassword string `json:"bind_password" toml:"bind_password" yaml:"bind_password" env:"BIND_PASSWORD"`
	// TLS是是否校验证书有效性
	SkipVerify bool `json:"skip_verify" toml:"skip_verify" yaml:"skip_verify" env:"SKIP_VERIFY"`
	// LDAP 服务器的登录用户名，必须是从根结点到用户节点的全路径
	BaseDn string `json:"base_dn" toml:"base_dn" yaml:"base_dn" env:"BASE_DN"`
	// 用户过滤条件
	UserFilter string `json:"user_filter" toml:"user_filter" yaml:"user_filter" env:"USER_FILTER"`
	// 用户组过滤条件
	GroupFilter string `json:"group_filter" toml:"group_filter" yaml:"group_filter" env:"GROUP_FILTER"`
	// 组属性的名称
	GroupNameAttribute string `json:"group_name_attribute" toml:"group_name_attribute" yaml:"group_name_attribute" env:"GROUP_NAME_ATTRIBUTE"`
	// 用户属性的名称
	UserNameAttribute string `json:"user_name_attribute" toml:"user_name_attribute" yaml:"user_name_attribute" env:"USER_NAME_ATTRIBUTE"`
	// 用户邮箱属性的名称
	MailAttribute string `json:"mail_attribute" toml:"mail_attribute" yaml:"mail_attribute" env:"MAIL_ATTRIBUTE"`
	// 用户显示名称属性名称
	DisplayNameAttribute string `json:"display_name_attribute" toml:"display_name_attribute" yaml:"display_name_attribute" env:"DISPLAY_NAME_ATTRIBUTE"`
	// 新增用户或者注销用户时，是否同步, 默认不做同步, 只读区用户信息
	SyncUser bool `json:"sync_user" toml:"sync_user" yaml:"sync_user" env:"SYNC_USER"`
}

func (c *Config) String() string {
	return pretty.ToJSON(c)
}
