package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
)

// 查询API接口列表
func (i *EndpointServiceImpl) QueryEndpoint(ctx context.Context, in *endpoint.QueryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	return nil, nil
}

// 查询API接口详情
func (i *EndpointServiceImpl) DescribeEndpoint(ctx context.Context, in *endpoint.DescribeEndpointRequest) (*endpoint.Endpoint, error) {
	return nil, nil
}

// 注册API接口
func (i *EndpointServiceImpl) RegistryEndpoint(ctx context.Context, in *endpoint.RegistryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	return nil, nil
}
