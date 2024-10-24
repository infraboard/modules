package endpoint

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
)

const (
	AppName = "endpoint"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	// 查询API接口列表
	QueryEndpoint(context.Context, *QueryEndpointRequest) (*types.Set[*Endpoint], error)
	// 查询API接口详情
	DescribeEndpoint(context.Context, *DescribeEndpointRequest) (*Endpoint, error)
	// 注册API接口
	RegistryEndpoint(context.Context, *RegistryEndpointRequest) (*types.Set[*Endpoint], error)
}

type QueryEndpointRequest struct {
}

type DescribeEndpointRequest struct {
}

type RegistryEndpointRequest struct {
	Items []*Endpoint `json:"items"`
}
