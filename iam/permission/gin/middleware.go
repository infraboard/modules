package permission

import (
	"context"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/gin/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/permission"
)

func Auth() gin.HandlerFunc {
	return NewAuther().Auth
}

func NewAuther() *Auther {
	return &Auther{
		tk: ioc.Controller().Get(token.AppName).(token.Service),
	}
}

type Auther struct {
	tk token.Service
}

func (a *Auther) Auth(c *gin.Context) {
	v := token.GetAccessTokenFromHTTP(c.Request)
	if v == "" {
		response.Failed(c, permission.ErrUnauthorized)
		c.Abort()
		return
	}

	in := token.NewValiateTokenRequest(v)
	tk, err := a.tk.ValiateToken(c.Request.Context(), in)
	if err != nil {
		response.Failed(c, err)
		c.Abort()
		return
	}

	ctx := context.WithValue(c.Request.Context(), token.CTX_TOKEN_KEY, tk)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// 写带参数的 Gin中间件
func Required(roles ...string) gin.HandlerFunc {
	p := &Permissoner{
		roles: roles,
	}
	return p.CheckPerm
}

type Permissoner struct {
	roles []string
}

func (p *Permissoner) CheckPerm(c *gin.Context) {
	tk := token.GetTokenFromCtx(c.Request.Context())
	if tk == nil {
		response.Failed(c, permission.ErrUnauthorized)
		return
	}

	if tk.IsAdmin {
		return
	}

	// err := p.HasPerm(tk.Role.String())
	// if err != nil {
	// 	response.Failed(c, err)
	// 	return
	// }
}

func (a *Permissoner) HasPerm(role string) error {
	if len(a.roles) == 0 {
		return nil
	}

	if slices.Contains(a.roles, role) {
		return nil
	}

	return exception.NewPermissionDeny("role %s not allow ", role)
}
