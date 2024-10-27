package issuer

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
)

func init() {
	ioc.Config().Registry(&PrivateTokenIssuer{})
}

type PrivateTokenIssuer struct {
	ioc.ObjectImpl

	user  user.Service
	token token.Service
}

func (p *PrivateTokenIssuer) Name() string {
	return "private_token_issuer"
}

func (p *PrivateTokenIssuer) Init() error {
	p.user = user.GetService()
	p.token = token.GetService()

	token.RegistryIssuer(token.ISSUER_PRIVATE_TOKEN, p)
	return nil
}

func (p *PrivateTokenIssuer) IssueToken(ctx context.Context, parameter token.IssueParameter) (*token.Token, error) {
	// 1. 查询用户
	uReq := user.NewDescribeUserRequestByUsername(parameter.Username())
	u, err := p.user.DescribeUser(ctx, uReq)
	if err != nil {
		if exception.IsNotFoundError(err) {
			return nil, exception.NewUnauthorized("%s", err)
		}
		return nil, err
	}

	if !u.EnabledApi {
		return nil, exception.NewPermissionDeny("未开启接口登录")
	}

	// 2. 校验Token合法
	_, err = p.token.ValiateToken(ctx, token.NewValiateTokenRequest(parameter.AccessToken()))
	if err != nil {
		return nil, err
	}

	// 3. 颁发token
	tk := token.NewToken()
	tk.UserId = u.Id
	tk.UserName = u.UserName
	tk.IsAdmin = u.IsAdmin

	expiredTTL := parameter.ExpireTTL()
	if expiredTTL > 0 {
		tk.SetExpiredAtByDuration(expiredTTL, 4)
	}
	return tk, nil
}
