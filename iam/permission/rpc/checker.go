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
	"github.com/rs/zerolog"
	"resty.dev/v3"
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

func GetPermissionChecker() *Checker {
	return ioc.Config().Get("permission_checker").(*Checker)
}

type Checker struct {
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (c *Checker) Name() string {
	return "permission_checker"
}

func (c *Checker) Priority() int {
	return gorestful.Priority() - 1
}

func (c *Checker) Init() error {
	c.log = log.Sub(c.Name())

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
	}

	next.ProcessFilter(r, w)
}

func (c *Checker) CheckToken(r *restful.Request) (*token.Token, error) {
	v := token.GetAccessTokenFromHTTP(r.Request)
	if v == "" {
		return nil, exception.NewUnauthorized("请先登录")
	}

	tk, err := c.ValiateToken(r.Request.Context(), token.NewValiateTokenRequest(v))
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Request.Context(), token.CTX_TOKEN_KEY, tk)
	r.Request = r.Request.WithContext(ctx)
	return tk, nil
}

// http://127.0.0.1:8020/api/mcenter/v1/token/validate
func (c *Checker) ValiateToken(ctx context.Context, in *token.ValiateTokenRequest) (*token.Token, error) {
	tk := token.NewToken()
	resp, err := resty.New().
		SetBaseURL(application.Get().InternalAddress).
		SetAuthToken(application.Get().InternalToken).
		R().
		WithContext(ctx).
		SetContentType("application/json").
		SetBody(in).
		SetResult(tk).
		Post("/api/mcenter/v1/token/validate")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, exception.NewUnauthorized("[%d] token校验异常: %s", resp.StatusCode(), resp.String())
	}
	if tk.NamespaceId == 0 {
		tk.NamespaceId = 1
	}
	return tk, nil
}

func (c *Checker) CheckPolicy(r *restful.Request, tk *token.Token, route *endpoint.RouteEntry) error {
	if tk.IsAdmin {
		return nil
	}

	// API权限校验
	if route.RequiredPerm {
		permReq := policy.NewValidateEndpointPermissionRequest()
		permReq.UserId = tk.UserId
		permReq.NamespaceId = tk.NamespaceId
		permReq.Service = application.Get().AppName
		permReq.Method = route.Method
		permReq.Path = route.Path
		resp, err := c.ValidateEndpointPermission(r.Request.Context(), permReq)
		if err != nil {
			return err
		}
		if !resp.HasPermission {
			return exception.NewPermissionDeny("无权限")
		}
	}

	return nil
}

// 查询策略列表
// /api/mcenter/v1/permission/check
func (c *Checker) ValidateEndpointPermission(ctx context.Context, in *policy.ValidateEndpointPermissionRequest) (*policy.ValidateEndpointPermissionResponse, error) {
	ins := policy.NewValidateEndpointPermissionResponse(*in)
	resp, err := resty.New().
		SetBaseURL(application.Get().InternalAddress).
		SetAuthToken(application.Get().InternalToken).
		SetDebug(true).
		R().
		WithContext(ctx).
		SetBody(in).
		SetResult(ins).
		Post("/api/mcenter/v1/permission/check")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode()/100 != 2 {
		return nil, exception.NewPermissionDeny("[%d] token鉴权异常: %s", resp.StatusCode(), resp.String())
	}
	return ins, nil
}
