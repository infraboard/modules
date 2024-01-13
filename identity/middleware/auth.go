package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_http "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/modules/identity/apps/token"
	"github.com/infraboard/modules/identity/apps/user"
)

func SetupAppHook() {
	// HTTP业务路由加载之前
	ioc_http.Get().GetRouterBuilder().BeforeLoadHooks(func(r http.Handler) {
		// GoRestful 框架
		if router, ok := r.(*gin.Engine); ok {
			// Gin Engine对象
			router.Use(NewTokenAuther().Auth)
		}
	})
}

func NewTokenAuther() *TokenAuther {
	return &TokenAuther{
		tk: ioc.Controller().Get(token.AppName).(token.Service),
	}
}

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type TokenAuther struct {
	tk   token.Service
	role user.Role
}

// 怎么鉴权?
// Gin中间件 func(*Context)
func (a *TokenAuther) Auth(c *gin.Context) {
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
}

// 权限鉴定, 鉴权是在用户已经认证的情况之下进行的
// 判断当前用户的角色
func (a *TokenAuther) Perm(c *gin.Context) {
	tkObj := c.Keys[token.TOKEN_GIN_KEY_NAME]
	if tkObj == nil {
		response.Failed(c.Writer, exception.NewPermissionDeny("token not found"))
		return
	}

	tk, ok := tkObj.(*token.Token)
	if !ok {
		response.Failed(c.Writer, exception.NewPermissionDeny("token not an *token.Token"))
		return
	}

	fmt.Printf("user %s role %d \n", tk.UserName, tk.Role)

	// 如果是Admin则直接放行
	if tk.Role == user.ROLE_ADMIN {
		return
	}

	if tk.Role != a.role {
		response.Failed(c.Writer, exception.NewPermissionDeny("role %d not allow", tk.Role))
		return
	}
}

// 写带参数的 Gin中间件
func Required(r user.Role) gin.HandlerFunc {
	a := NewTokenAuther()
	a.role = r
	return a.Perm
}
