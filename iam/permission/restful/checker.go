package permission

import (
	"context"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/rs/zerolog"
)

func Auth(v bool) (string, bool) {
	return endpoint.META_REQUIRED_AUTH_KEY, v
}

func Permission(v bool) (string, bool) {
	return endpoint.META_REQUIRED_PERM_KEY, v
}

func Required(roles ...string) (string, []string) {
	return endpoint.META_REQUIRED_ROLE_KEY, roles
}

func init() {
	ioc.Config().Registry(&Checker{})
}

type Checker struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	token  token.Service
	policy policy.Service
}

func (c *Checker) Name() string {
	return "permission_checker"
}

func (c *Checker) Priority() int {
	return gorestful.Priority() - 1
}

func (c *Checker) Init() error {
	c.log = log.Sub(c.Name())
	c.token = token.GetService()
	c.policy = policy.GetService()

	// 注册认证中间件
	gorestful.RootRouter().Filter(c.Check)
	return nil
}

func (c *Checker) Check(r *restful.Request, w *restful.Response, next *restful.FilterChain) {
	route := endpoint.NewEntryFromRestRoute(r.SelectedRoute())
	if route.RequiredAuth {
		// 校验身份
		tk, err := c.CheckToken(r)
		if err != nil {
			response.Failed(w, err)
			return
		}

		// 校验权限
		if err := c.CheckPolicy(r, tk, route); err != nil {
			response.Failed(w, err)
			return
		}
	}

	next.ProcessFilter(r, w)
}

func (c *Checker) CheckToken(r *restful.Request) (*token.Token, error) {
	v := token.GetAccessTokenFromHTTP(r.Request)
	if v == "" {
		return nil, exception.NewUnauthorized("请先登录")
	}

	tk, err := c.token.ValiateToken(r.Request.Context(), token.NewValiateTokenRequest(v))
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Request.Context(), token.CTX_TOKEN_KEY, tk)
	r.Request = r.Request.WithContext(ctx)
	return tk, nil
}

func (c *Checker) CheckPolicy(r *restful.Request, tk *token.Token, route *endpoint.RouteEntry) error {
	if tk.IsAdmin {
		return nil
	}

	set, err := c.policy.QueryPolicy(r.Request.Context(),
		policy.NewQueryPolicyRequest().
			SetNamespaceId(tk.NamespaceId).
			SetUserId(tk.UserId).
			SetExpired(false).
			SetEnabled(true).
			SetWithRole(true),
	)
	if err != nil {
		return exception.NewInternalServerError(err.Error())
	}

	// 角色校验
	if route.HasRequiredRole() {
		for i := range set.Items {
			p := set.Items[i]
			if route.IsRequireRole(p.Role.Name) {
				return nil
			}
		}
	}

	// API权限校验
	if route.RequiredPerm {
		for i := range set.Items {
			p := set.Items[i]
			if p.Role == nil {
				return exception.NewInternalServerError("policy role is nil")
			}

		}
	}

	return nil
}