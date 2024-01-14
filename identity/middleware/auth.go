package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/identity/apps/token"
	"github.com/infraboard/modules/identity/apps/user"
)

func NewTokenAuther(auth bool, roles []user.Role) *TokenAuther {
	return &TokenAuther{
		tk:    ioc.Controller().Get(token.AppName).(token.Service),
		auth:  auth,
		roles: roles,
	}
}

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type TokenAuther struct {
	tk    token.Service
	auth  bool
	roles []user.Role
}

// 怎么鉴权?
// Gin中间件 func(*Context)
func (a *TokenAuther) Auth(c *gin.Context) {
	if !a.auth {
		return
	}

	// 1. 获取Token
	at, err := c.Cookie(token.TOKEN_COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			response.Failed(c.Writer, token.CookieNotFound)
			return
		}
		response.Failed(c.Writer, err)
		return
	}

	// 2.调用Token模块来认证
	in := token.NewValiateToken(at)
	tk, err := a.tk.ValiateToken(c.Request.Context(), in)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 把鉴权后的 结果: tk, 放到请求的上下文, 方便后面的业务逻辑使用
	if c.Keys == nil {
		c.Keys = map[string]any{}
	}
	c.Keys[token.TOKEN_GIN_KEY_NAME] = tk

	// 权限鉴定, 鉴权是在用户已经认证的情况之下进行的
	// 判断当前用户的角色
	if tk.Role == user.ROLE_ADMIN || len(a.roles) == 0 {
		return
	}

	err = a.HasPerm(tk.Role.String())
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
}

func (a *TokenAuther) HasPerm(role string) error {
	if len(a.roles) == 0 {
		return nil
	}

	for _, r := range a.roles {
		if r.String() == role {
			return nil
		}
	}

	return exception.NewPermissionDeny("role %s not allow ", role)
}

// 写带参数的 Gin中间件
func Required(auth bool, roles ...user.Role) gin.HandlerFunc {
	a := NewTokenAuther(auth, roles)
	return a.Auth
}
