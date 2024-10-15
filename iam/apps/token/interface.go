package token

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "tokens"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 登录接口(颁发Token)
	IssueToken(context.Context, *IssueTokenRequest) (*Token, error)
	// 退出接口(销毁Token)
	RevolkToken(context.Context, *RevolkTokenRequest) (*Token, error)

	// 校验Token 是给内部中间层使用 身份校验层
	// 校验完后返回Token, 通过Token获取 用户信息
	ValiateToken(context.Context, *ValiateToken) (*Token, error)
}

func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{
		Parameter: make(ProviderParameter),
	}
}

type IssueTokenRequest struct {
	// 认证方式
	Provider string `json:"provider"`
	// 参数
	Parameter ProviderParameter `json:"parameter"`
}

type ProviderParameter map[string]string

func (p ProviderParameter) Username() string {
	return p["username"]
}

func (p ProviderParameter) Password() string {
	return p["password"]
}

func NewRevolkTokenRequest(at, rk string) *RevolkTokenRequest {
	return &RevolkTokenRequest{
		AccessToken:  at,
		RefreshToken: rk,
	}
}

// 万一的Token泄露, 不知道refresh_token，也没法推出
type RevolkTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewValiateToken(at string) *ValiateToken {
	return &ValiateToken{
		AccessToken: at,
	}
}

type ValiateToken struct {
	AccessToken string `json:"access_token"`
}
