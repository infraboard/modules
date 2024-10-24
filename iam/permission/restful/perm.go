package permission

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/policy"
	"github.com/infraboard/modules/iam/apps/token"
)

func Permission(v bool) (string, bool) {
	return endpoint.META_REQUIRED_PERM_KEY, v
}

func Required(roles ...string) (string, []string) {
	return endpoint.META_REQUIRED_ROLE_KEY, roles
}

func NewPermissoner() *Permissoner {
	return &Permissoner{
		policy: policy.GetService(),
	}
}

type Permissoner struct {
	policy policy.Service
}

func (p *Permissoner) CheckPerm(r *restful.Request, w *restful.Response, next *restful.FilterChain) {
	tk := token.GetTokenFromCtx(r.Request.Context())
	if tk == nil {
		response.Failed(w, exception.NewUnauthorized("请先登录"))
		return
	}

	if tk.IsAdmin {
		return
	}

	set, err := p.policy.QueryPolicy(r.Request.Context(),
		policy.NewQueryPolicyRequest().
			SetNamespaceId(tk.NamespaceId).
			SetUserId(tk.UserId).
			SetExpired(false).
			SetEnabled(true).
			SetWithRole(true),
	)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError(err.Error()))
		return
	}

	var hasPerm bool
	route := endpoint.NewEntryFromRestRoute(r.SelectedRoute())
	if route.HasRequiredRole() {
		set.ForEach(func(t *policy.Policy) {
			if route.IsRequireRole(t.Role.Name) {
				hasPerm = true
			}
		})
		if hasPerm {
			response.Success(w, set)
			return
		}
	}
}
