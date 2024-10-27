package issuer

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
)

func init() {
	ioc.Config().Registry(&PasswordTokenIssuer{
		ExpiredTTLSecond: 1 * 60 * 60,
	})
}

type PasswordTokenIssuer struct {
	ioc.ObjectImpl
	user user.Service

	// Password颁发的Token 过去时间由系统配置, 不允许用户自己设置
	ExpiredTTLSecond int `json:"expired_ttl_second" toml:"expired_ttl_second" yaml:"expired_ttl_second" env:"EXPIRED_TTL_SECOND"`

	expiredDuration time.Duration
}

func (p *PasswordTokenIssuer) Name() string {
	return "password_token_issuer"
}

func (p *PasswordTokenIssuer) Init() error {
	p.user = user.GetService()
	p.expiredDuration = time.Duration(p.ExpiredTTLSecond) * time.Second

	token.RegistryIssuer(token.ISSUER_PASSWORD, p)
	return nil
}

func (p *PasswordTokenIssuer) IssueToken(ctx context.Context, parameter token.IssueParameter) (*token.Token, error) {
	// 1. 查询用户
	uReq := user.NewDescribeUserRequestByUsername(parameter.Username())
	u, err := p.user.DescribeUser(ctx, uReq)
	if err != nil {
		if exception.IsNotFoundError(err) {
			return nil, exception.NewUnauthorized("%s", err)
		}
		return nil, err
	}

	// 2. 比对密码
	err = u.CheckPassword(parameter.Password())
	if err != nil {
		return nil, exception.NewUnauthorized("%s", err)
	}

	// 3. 颁发token
	tk := token.NewToken()
	tk.UserId = u.Id
	tk.UserName = u.UserName
	tk.IsAdmin = u.IsAdmin

	tk.SetExpiredAtByDuration(p.expiredDuration, 4)
	return tk, nil
}
