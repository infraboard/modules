package permission

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"resty.dev/v3"
)

func init() {
	ioc.Api().Registry(&ApiRegister{})
}

func GetApiRegister() *ApiRegister {
	return ioc.Api().Get("api_register").(*ApiRegister)
}

type ApiRegister struct {
	ioc.ObjectImpl
}

func (c *ApiRegister) Name() string {
	return "api_register"
}

func (i *ApiRegister) Priority() int {
	return -100
}

func (a *ApiRegister) Init() error {
	// 注册认证中间件
	entries := endpoint.NewEntryFromRestfulContainer(gorestful.RootRouter())
	req := endpoint.NewRegistryEndpointRequest()
	req.AddItem(entries...)
	_, err := a.RegistryEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

// 注册API接口(RPC --> REST SDK)
// 自己的 注册API接口
// restful client: github.com/go-resty/resty/v2
// http://127.0.0.1:8020/api/mcenter/v1/endpoint
func (a *ApiRegister) RegistryEndpoint(ctx context.Context, in *endpoint.RegistryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	set := types.New[*endpoint.Endpoint]()
	resp, err := resty.New().
		SetDebug(true).
		SetBaseURL(application.Get().InternalAddress).
		SetAuthToken(application.Get().InternalToken).
		R().
		WithContext(ctx).
		SetContentType("application/json").
		SetBody(in.Items).
		SetResult(set).
		Post("/api/mcenter/v1/endpoint")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode()/100 != 2 {
		return nil, exception.NewPermissionDeny("[%d] API注册异常: %s", resp.StatusCode(), resp.String())
	}
	return set, nil
}
