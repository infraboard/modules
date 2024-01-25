package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/infraboard/modules/identity/apps/user"
	"github.com/infraboard/modules/identity/middleware"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	// 引入模块
	_ "github.com/infraboard/modules/identity"
	_ "github.com/infraboard/modules/identity/cmd"
)

func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&ApiHandler{})

	server.Root.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "example API服务",
			Run: func(cmd *cobra.Command, args []string) {
				cobra.CheckErr(server.Run(context.Background()))
			},
		},
	)

	// 启动
	server.RunCLI()
}

type ApiHandler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl

	// mysql db依赖
	db *gorm.DB
}

// 覆写对象的名称, 该名称名称会体现在API的路径前缀里面
// 比如: /simple/api/v1/module_a/db_stats
// 其中/simple/api/v1/module_a 就是对象API前缀, 命名规则如下:
// <service_name>/<path_prefix>/<object_version>/<object_name>
func (h *ApiHandler) Name() string {
	return "module_a"
}

// 初始化db属性, 从ioc的配置区域获取共用工具 gorm db对象
func (h *ApiHandler) Init() error {
	h.db = datasource.DB()
	return nil
}

// API路由
func (h *ApiHandler) Registry(r gin.IRouter) {
	r.Use(middleware.Auth())
	r.GET("/db_stats", middleware.Perm(user.ROLE_MEMBER), h.DBStats)
}

func (h *ApiHandler) DBStats(ctx *gin.Context) {
	db, _ := h.db.DB()
	ctx.JSON(http.StatusOK, gin.H{
		"data": db.Stats(),
	})
}
