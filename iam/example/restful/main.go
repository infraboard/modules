package main

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
	"gorm.io/gorm"

	// 引入模块
	_ "github.com/infraboard/modules/iam/init"

	// 非功能性模块
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/restful"
)

func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&ApiHandler{})

	// 启动
	cmd.Start()
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

	tags := []string{"DB状态"}
	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("/db_stats").To(h.DBStats).
		Doc("查询数据库状态").
		Metadata(restfulspec.KeyOpenAPITags, tags))

	return nil
}

func (h *ApiHandler) DBStats(r *restful.Request, w *restful.Response) {
	db, _ := h.db.DB()
	w.WriteAsJson(gin.H{
		"data": db.Stats(),
	})
}
