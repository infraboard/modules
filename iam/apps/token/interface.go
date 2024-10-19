package token

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "token"
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
	ValiateToken(context.Context, *ValiateTokenRequest) (*Token, error)

	// 查询已经颁发出去的Token
	QueryToken(context.Context, *QueryTokenRequest) (*types.Set[*Token], error)
}

func NewQueryTokenRequest() *QueryTokenRequest {
	return &QueryTokenRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryTokenRequest struct {
	*request.PageRequest
}

func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{
		Parameter: make(IssueParameter),
	}
}

type IssueTokenRequest struct {
	// 认证方式
	Issuer string `json:"issuer"`
	// 参数
	Parameter IssueParameter `json:"parameter"`
}

func (i *IssueTokenRequest) IssueByPassword(username, password string) {
	i.Issuer = ISSUER_PASSWORD
	i.Parameter.SetUsername(username)
	i.Parameter.SetPassword(password)
}

func GetIssueParameterValue[T any](p IssueParameter, key string) T {
	types.New[*Token]()
	v := p[key]
	if v != nil {
		if value, ok := v.(T); ok {
			return value
		}
	}
	var zero T
	return zero
}

type IssueParameter map[string]any

/*
password issuer parameter
*/

func (p IssueParameter) Username() string {
	return GetIssueParameterValue[string](p, "username")
}

func (p IssueParameter) Password() string {
	return GetIssueParameterValue[string](p, "password")
}

func (p IssueParameter) SetUsername(v string) {
	p["username"] = v
}

func (p IssueParameter) SetPassword(v string) {
	p["password"] = v
}

/*
private token issuer parameter
*/

func (p IssueParameter) AccessToken() string {
	return GetIssueParameterValue[string](p, "access_token")
}

func (p IssueParameter) ExpireTTL() time.Duration {
	return time.Second * time.Duration(GetIssueParameterValue[int64](p, "expired_ttl"))
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

func NewValiateTokenRequest(accessToken string) *ValiateTokenRequest {
	return &ValiateTokenRequest{
		AccessToken: accessToken,
	}
}

type ValiateTokenRequest struct {
	AccessToken string `json:"access_token"`
}
