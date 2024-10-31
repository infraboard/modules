package ldap

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&LdapTokenIssuer{
		ExpiredTTLSecond: 1 * 60 * 60,
	})
}

type LdapTokenIssuer struct {
	ioc.ObjectImpl
	log  *zerolog.Logger
	user user.Service

	// Password颁发的Token 过去时间由系统配置, 不允许用户自己设置
	ExpiredTTLSecond int `json:"expired_ttl_second" toml:"expired_ttl_second" yaml:"expired_ttl_second" env:"EXPIRED_TTL_SECOND"`
	// Ldap
	Config

	expiredDuration time.Duration
}

func (p *LdapTokenIssuer) Name() string {
	return "ldap_token_issuer"
}

func (p *LdapTokenIssuer) Init() error {
	p.log = log.Sub(p.Name())
	p.user = user.GetService()
	p.expiredDuration = time.Duration(p.ExpiredTTLSecond) * time.Second

	token.RegistryIssuer(token.ISSUER_LDAP, p)
	return nil
}

func (i *LdapTokenIssuer) IssueToken(ctx context.Context, parameter token.IssueParameter) (*token.Token, error) {
	// 连接Ldap Server
	p := NewLdapProvider(i.Config)
	err := p.CheckConnect()
	if err != nil {
		return nil, err
	}

	// 检查用户密码是否正确
	u, err := p.CheckUserPassword(parameter.Username(), parameter.Password())
	if err != nil {
		return nil, err
	}

	// 判断用户是否在数据库存在, 如果不存在需要同步到本地数据库
	lu, err := i.user.DescribeUser(ctx, user.NewDescribeUserRequestByUserName(u.Username))
	if err != nil {
		if exception.IsNotFoundError(err) {
			i.log.Debug().Msgf("sync user: (%s) to db", u.Username)
			// 创建本地用户, 密码等同于LDAP密码
			newReq := user.NewLDAPCreateUserRequest(u.Username, parameter.Password(), "系统自动生成")
			lu, err = i.user.CreateUser(ctx, newReq)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// 3. 颁发token
	tk := token.NewToken()
	tk.UserId = lu.Id
	tk.UserName = lu.UserName
	tk.IsAdmin = lu.IsAdmin

	tk.SetExpiredAtByDuration(i.expiredDuration, 4)
	return tk, nil
}
