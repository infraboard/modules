package provider

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
)

func init() {
	ioc.Config().Registry(&PasswordTokenProvider{})
}

type PasswordTokenProvider struct {
	ioc.ObjectImpl

	user user.Service
}

func (p *PasswordTokenProvider) Name() string {
	return "password_token_provider"
}

func (p *PasswordTokenProvider) Init() error {
	p.user = user.GetService()

	token.RegistryProvider("password", p)
	return nil
}

func (p *PasswordTokenProvider) IssueToken(ctx context.Context, parameter token.ProviderParameter) (*token.Token, error) {
	// 1. 查询用户
	uReq := user.NewDescribeUserRequestByUsername(parameter.Username())
	u, err := p.user.DescribeUser(ctx, uReq)
	if err != nil {
		if exception.IsNotFoundError(err) {
			return nil, token.NewAuthFailed("%s", err)
		}
		return nil, err
	}

	// 2. 比对密码
	err = u.CheckPassword(parameter.Password())
	if err != nil {
		return nil, token.NewAuthFailed("%s", err)
	}

	// 3. 颁发token
	tk := token.NewToken()
	tk.UserId = fmt.Sprintf("%d", u.Id)
	tk.UserName = u.Username
	tk.IsAdmin = u.IsAdmin
	return tk, nil
}
