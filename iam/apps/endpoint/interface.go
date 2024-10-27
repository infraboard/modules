package endpoint

import (
	"context"

	"github.com/infraboard/mcube/v2/http/request"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps"
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

func NewQueryEndpointRequest() *QueryEndpointRequest {
	return &QueryEndpointRequest{
		PageRequest: request.NewDefaultPageRequest(),
	}
}

type QueryEndpointRequest struct {
	*request.PageRequest
}

func NewDescribeEndpointRequest() *DescribeEndpointRequest {
	return &DescribeEndpointRequest{}
}

type DescribeEndpointRequest struct {
	apps.GetRequest
}

func NewRegistryEndpointRequest() *RegistryEndpointRequest {
	return &RegistryEndpointRequest{
		Items: []*Endpoint{},
	}
}

type RegistryEndpointRequest struct {
	Items []*Endpoint `json:"items"`
}
