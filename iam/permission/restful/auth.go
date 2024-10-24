package permission

import (
	"context"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/permission"
)

func Auth(v bool) (string, bool) {
	return endpoint.META_REQUIRED_AUTH_KEY, v
}

func NewAuther() *Auther {
	return &Auther{
		tk: ioc.Controller().Get(token.AppName).(token.Service),
	}
}

type Auther struct {
	tk token.Service
}

func (a *Auther) Auth(r *restful.Request, w *restful.Response, next *restful.FilterChain) {
	v := token.GetAccessTokenFromHTTP(r.Request)
	if v == "" {
		response.Failed(w, permission.ErrUnauthorized)
		return
	}

	// 2.调用Token模块来认证
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
