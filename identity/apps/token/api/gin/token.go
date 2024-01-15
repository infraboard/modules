package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/modules/identity/apps/token"
)

// 需要把HandleFunc 添加到Root路由，定义 API ---> HandleFunc
// 可以选择把这个Handler上的HandleFunc都注册到路由上面
func (h *TokenApiHandler) Registry(r gin.IRouter) {
	// r 是Gin的路由器
	r.POST("/", h.Login)
	r.DELETE("/", h.Logout)
}

func (h *TokenApiHandler) Login(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := token.NewLoginRequest()

	err := c.BindJSON(req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	c.SetCookie(token.ACCESS_TOKEN_COOKIE_NAME, tk.AccessToken, 0, "/", application.Get().Domain, false, true)

	// 3. 返回响应
	response.Success(c.Writer, tk)
}

// Logout HandleFunc
func (h *TokenApiHandler) Logout(c *gin.Context) {
	req := token.NewLogoutRequest(
		token.GetAccessTokenFromHTTP(c.Request),
		token.GetRefreshTokenFromHTTP(c.Request),
	)

	tk, err := h.svc.Logout(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	c.SetCookie(token.ACCESS_TOKEN_COOKIE_NAME, "", 0, "/", application.Get().Domain, false, true)

	// 3. 返回响应
	response.Success(c.Writer, tk)
}
