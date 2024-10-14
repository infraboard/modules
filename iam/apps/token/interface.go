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
	Login(context.Context, *LoginRequest) (*Token, error)
	// 退出接口(销毁Token)
	Logout(context.Context, *LogoutRequest) (*Token, error)

	// 校验Token 是给内部中间层使用 身份校验层
	// 校验完后返回Token, 通过Token获取 用户信息
	ValiateToken(context.Context, *ValiateToken) (*Token, error)
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

type LoginRequest struct {
	Username string
	Password string
}

func NewLogoutRequest(at, rk string) *LogoutRequest {
	return &LogoutRequest{
		AccessToken:  at,
		RefreshToken: rk,
	}
}

// 万一的Token泄露, 不知道refresh_token，也没法推出
type LogoutRequest struct {
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
