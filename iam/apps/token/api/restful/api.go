package restful

import (
	"github.com/infraboard/mcube/v2/http/label"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/modules/iam/apps/token"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
)

func init() {
	ioc.Api().Registry(&TokenRestulApiHandler{})
}

type TokenRestulApiHandler struct {
	ioc.ObjectImpl

	// 依赖控制器
	svc token.Service
}

func (h *TokenRestulApiHandler) Name() string {
	return token.AppName
}

func (h *TokenRestulApiHandler) Init() error {
	h.svc = token.GetService()

	tags := []string{"登录"}
	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.POST("").To(h.Login).
		Doc("颁发令牌").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(token.IssueTokenRequest{}).
		Writes(token.Token{}).
		Returns(200, "OK", token.Token{}))

	ws.Route(ws.DELETE("").To(h.Logout).
		Doc("撤销令牌").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Auth, label.Enable).
		Metadata(label.PERMISSION_MODE, label.PERMISSION_MODE_ACL.Value()).
		Reads(token.IssueTokenRequest{}).
		Writes(token.Token{}).
		Returns(200, "OK", token.Token{}).
		Returns(404, "Not Found", nil))
	return nil
}
