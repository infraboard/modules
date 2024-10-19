package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/modules/iam/apps/token"
)

func (h *TokenGinApiHandler) Login(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := token.NewIssueTokenRequest()

	err := c.BindJSON(req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.IssueToken(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	c.SetCookie(token.ACCESS_TOKEN_COOKIE_NAME, tk.AccessToken,
		tk.AccessTokenExpiredTTL(), "/", application.Get().Domain(), false, true)
	// 在Header头中也添加Token
	c.Header(token.ACCESS_TOKEN_RESPONSE_HEADER_NAME, tk.AccessToken)

	// 3. Body中返回Token对象
	response.Success(c.Writer, tk)
}

// Logout HandleFunc
func (h *TokenGinApiHandler) Logout(c *gin.Context) {
	req := token.NewRevolkTokenRequest(
		token.GetAccessTokenFromHTTP(c.Request),
		token.GetRefreshTokenFromHTTP(c.Request),
	)

	tk, err := h.svc.RevolkToken(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// access_token 通过SetCookie 直接写到浏览器客户端(Web)
	c.SetCookie(token.ACCESS_TOKEN_COOKIE_NAME, "", 0, "/", application.Get().Domain(), false, true)

	// 3. 返回响应
	response.Success(c.Writer, tk)
}
