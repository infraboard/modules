package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/modules/identity/apps/token"
)

// 需要把HandleFunc 添加到Root路由，定义 API ---> HandleFunc
// 可以选择把这个Handler上的HandleFunc都注册到路由上面
func (h *TokenApiHandler) Registry(r gin.IRouter) {
	// r 是Gin的路由器
	v1 := r.Group("v1")
	v1.POST("/tokens/", h.Login)
	v1.DELETE("/tokens/", h.Logout)
}

// Login HandleFunc
func (h *TokenApiHandler) Login(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	// 一定要使用JSON
	req := token.NewLoginRequest()

	// json.unmarsal
	// http boyd ---> LoginRequest Object
	err := c.BindJSON(req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 2. 执行逻辑
	// 把http 协议的请求 ---> 控制器的请求
	ins, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	c.SetCookie(token.TOKEN_COOKIE_NAME, ins.AccessToken, 0, "/", "localhost", false, true)

	// 3. 返回响应
	response.Success(c.Writer, ins)
}

// Logout HandleFunc
func (h *TokenApiHandler) Logout(*gin.Context) {

}
