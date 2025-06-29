package permission

import (
	"context"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/permission"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&Checker{})
}

func Auth(v bool) (string, bool) {
	return endpoint.META_REQUIRED_AUTH_KEY, v
}

func Permission(v bool) (string, bool) {
	return endpoint.META_REQUIRED_PERM_KEY, v
}

func Resource(v string) (string, string) {
	return endpoint.META_RESOURCE_KEY, v
}

func Action(v string) (string, string) {
	return endpoint.META_ACTION_KEY, v
}

func Required(roles ...string) (string, []string) {
	return endpoint.META_REQUIRED_ROLE_KEY, roles
}

type Checker struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	token  token.Service
	policy policy.Service
}

func (c *Checker) Name() string {
	return CHECKER_APP_NAME
}

func (c *Checker) Priority() int {
	return permission.GetCheckerPriority()
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
	route := endpoint.NewEntryFromRestRouteReader(r.SelectedRoute())
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

		// 添加上下文
		ctx := context.WithValue(r.Request.Context(), token.CTX_TOKEN_KEY, tk)
		r.Request = r.Request.WithContext(ctx)
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
	return tk, nil
}

func (c *Checker) CheckPolicy(r *restful.Request, tk *token.Token, route *endpoint.RouteEntry) error {
	if tk.IsAdmin {
		return nil
	}

	// 角色校验
	if route.HasRequiredRole() {
		set, err := c.policy.QueryPolicy(r.Request.Context(),
			policy.NewQueryPolicyRequest().
				SetNamespaceId(tk.GetNamespaceId()).
				SetUserId(tk.UserId).
				SetExpired(false).
				SetEnabled(true).
				SetWithRole(true),
		)
		if err != nil {
			return exception.NewInternalServerError("%s", err.Error())
		}
		hasPerm := false
		for i := range set.Items {
			p := set.Items[i]
			if route.IsRequireRole(p.Role.Name) {
				hasPerm = true
				break
			}
		}
		if !hasPerm {
			return exception.NewPermissionDeny("无权限访问")
		}
	}

	// API权限校验
	if route.RequiredPerm {
		validateReq := policy.NewValidateEndpointPermissionRequest()
		validateReq.UserId = tk.UserId
		validateReq.ResourceScope = tk.ResourceScope
		validateReq.Service = application.Get().GetAppName()
		validateReq.Method = route.Method
		validateReq.Path = route.Path
		resp, err := c.policy.ValidateEndpointPermission(r.Request.Context(), validateReq)
		if err != nil {
			return exception.NewInternalServerError("%s", err.Error())
		}
		if !resp.HasPermission {
			return exception.NewPermissionDeny("无权限访问")
		}

		tk.ResourceScope = resp.ResourceScope
		tk.BuildMySQLPrefixBlob()
	}

	return nil
}
