package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/modules/identity/apps/user"
	"github.com/infraboard/modules/identity/middleware"
)

func (h *UserApiHandler) Registry(r gin.IRouter) {
	r.Use(middleware.Auth(), middleware.Perm(user.ROLE_ADMIN))
	r.GET("/", h.QueryUser)
	r.GET("/:id", h.DescribeUser)
	r.POST("/", h.CreateUser)
	r.DELETE("/:id", h.DeleteUser)
}

func (h *UserApiHandler) QueryUser(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewQueryUserRequest()
	req.PageRequest = request.NewPageRequestFromHTTP(c.Request)

	// 2. 执行逻辑
	tk, err := h.svc.QueryUser(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 3. 返回响应
	response.Success(c.Writer, tk)
}

func (h *UserApiHandler) DescribeUser(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewDescribeUserRequestById(c.Param("id"))

	// 2. 执行逻辑
	tk, err := h.svc.DescribeUser(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 3. 返回响应
	response.Success(c.Writer, tk)
}

func (h *UserApiHandler) CreateUser(c *gin.Context) {
	// 1. 获取用户的请求参数， 参数在Body里面
	req := user.NewCreateUserRequest()

	err := c.BindJSON(req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 2. 执行逻辑
	tk, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 3. 返回响应
	response.Success(c.Writer, tk)
}

func (h *UserApiHandler) DeleteUser(c *gin.Context) {
	req := user.NewDeleteUserRequest(c.Param("id"))
	u, err := h.svc.DeleteUser(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	// 3. 返回响应
	response.Success(c.Writer, u)
}
