package feishu

import "net/url"

// NewDefaultConfig represents the default LDAP config.
func NewDefaultConfig() *Config {
	return &Config{}
}

type Config struct {
	// 开启飞书认证
	Enabled bool `json:"enabled" toml:"enabled" yaml:"enabled" env:"ENABLED"`
	// 飞书应用凭证, Oauth2.0时 也叫client_id
	AppId string `json:"app_id" toml:"app_id" yaml:"app_id" env:"API_ID"`
	// 飞书应用凭证, Oauth2.0时 也叫client_secret
	AppSecret string `json:"app_secret" toml:"app_secret" yaml:"app_secret" env:"API_SECRET"`
	// Oauth2.0时, 应用服务地址页面
	RedirectUri string `json:"redirect_uri" toml:"redirect_uri" yaml:"redirect_uri" env:"REDIRECT_URI"`
}

func (c *Config) MakeGetTokenFormRequest(code string) string {
	form := make(url.Values)
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", c.AppId)
	form.Add("client_secret", c.AppSecret)
	form.Add("code", code)
	form.Add("redirect_uri", c.RedirectUri)
	return form.Encode()
}
