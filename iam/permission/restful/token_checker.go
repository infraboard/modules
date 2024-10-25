package permission

import (
	"context"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/permission"
	"github.com/rs/zerolog"
)

func Auth(v bool) (string, bool) {
	return endpoint.META_REQUIRED_AUTH_KEY, v
}

func init() {
	ioc.Config().Registry(&TokenChecker{})
}

type TokenChecker struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	tk token.Service
}

func (p *TokenChecker) Name() string {
	return "token_checker"
}

func (m *TokenChecker) Priority() int {
	return gorestful.Priority() - 1
}

func (p *TokenChecker) Init() error {
	p.log = log.Sub(p.Name())
	p.tk = token.GetService()

	// 注册认证中间件
	gorestful.RootRouter().Filter(p.Auth)
	return nil
}

func (a *TokenChecker) Auth(r *restful.Request, w *restful.Response, next *restful.FilterChain) {
	v := token.GetAccessTokenFromHTTP(r.Request)
	if v == "" {
		response.Failed(w, permission.ErrUnauthorized)
		return
	}

	in := token.NewValiateTokenRequest(v)
	tk, err := a.tk.ValiateToken(r.Request.Context(), in)
	if err != nil {
		response.Failed(w, err)
		return
	}

	ctx := context.WithValue(r.Request.Context(), token.CTX_TOKEN_KEY, tk)
	r.Request = r.Request.WithContext(ctx)
	next.ProcessFilter(r, w)
}
