package permission

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/gin/response"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
	"github.com/infraboard/modules/iam/permission"
)

const (
	META_AUTH_KEY     = "auth"
	META_REQUIRED_KEY = "required"
)

func Auth(v bool) (string, bool) {
	return META_AUTH_KEY, v
}

// 写带参数的 Gin中间件
func Required(roles ...user.Role) (string, []user.Role) {
	return META_REQUIRED_KEY, roles
}

func NewAuther() *Auther {
	return &Auther{
		tk: ioc.Controller().Get(token.AppName).(token.Service),
	}
}

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type Auther struct {
	tk token.Service
}

// 怎么鉴权?
// Gin中间件 func(*Context)
func (a *Auther) Auth(c *gin.Context) {
	// 1. 获取Token
	v := token.GetAccessTokenFromHTTP(c.Request)
	if v == "" {
		response.Failed(c, permission.ErrUnauthorized)
		return
	}

	// 2.调用Token模块来认证
	in := token.NewValiateToken(v)
	tk, err := a.tk.ValiateToken(c.Request.Context(), in)
	if err != nil {
		response.Failed(c, err)
		return
	}

	// 把鉴权后的 结果: tk, 放到请求的上下文, 方便后面的业务逻辑使用
	c.Set(token.ACCESS_TOKEN_GIN_KEY_NAME, tk)
}

type Permissoner struct {
	roles []user.Role
}

func (p *Permissoner) CheckPerm(c *gin.Context) {
	v := c.Keys[token.ACCESS_TOKEN_GIN_KEY_NAME]
	if v == nil {
		response.Failed(c, permission.ErrUnauthorized)
		return
	}

	tk, ok := v.(*token.Token)
	if !ok {
		response.Failed(c, fmt.Errorf("tk not *token.Token"))
		return
	}

	// 权限鉴定, 鉴权是在用户已经认证的情况之下进行的
	// 判断当前用户的角色
	if tk.Role == user.ROLE_ADMIN {
		return
	}

	err := p.HasPerm(tk.Role.String())
	if err != nil {
		response.Failed(c, err)
		return
	}
}

func (a *Permissoner) HasPerm(role string) error {
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
