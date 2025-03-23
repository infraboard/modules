package permission

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&ApiRegister{})
}

func GetApiRegister() *ApiRegister {
	return ioc.Api().Get("api_register").(*ApiRegister)
}

type ApiRegister struct {
	ioc.ObjectImpl

	log *zerolog.Logger
}

func (c *ApiRegister) Name() string {
	return "api_register"
}

func (i *ApiRegister) Priority() int {
	return -100
}

func (a *ApiRegister) Init() error {
	a.log = log.Sub(a.Name())
	// 注册认证中间件
	entries := endpoint.NewEntryFromRestfulContainer(gorestful.RootRouter())
	req := endpoint.NewRegistryEndpointRequest()
	req.AddItem(entries...)
	set, err := endpoint.GetService().RegistryEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	a.log.Info().Msgf("registry endpoinst: %s", set.Items)
	return nil
}
