package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/endpoint"
	"gorm.io/gorm"
)

// 注册API接口
func (i *EndpointServiceImpl) RegistryEndpoint(ctx context.Context, in *endpoint.RegistryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	set := types.New[*endpoint.Endpoint]()
	err := datasource.DBFromCtx(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range in.Items {
			item := in.Items[i]
			ins := endpoint.NewEndpoint().SetRouteEntry(*item)

			if err := tx.Save(ins).Error; err != nil {
				return err
			}
			set.Add(ins)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return set, nil
}

// 查询API接口列表
func (i *EndpointServiceImpl) QueryEndpoint(ctx context.Context, in *endpoint.QueryEndpointRequest) (*types.Set[*endpoint.Endpoint], error) {
	set := types.New[*endpoint.Endpoint]()

	query := datasource.DBFromCtx(ctx).Model(&endpoint.Endpoint{})
	err := query.Count(&set.Total).Error
	if err != nil {
		return nil, err
	}

	err = query.
		Order("created_at desc").
		Offset(int(in.ComputeOffset())).
		Limit(int(in.PageSize)).
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}
	return set, nil
}

// 查询API接口详情
func (i *EndpointServiceImpl) DescribeEndpoint(ctx context.Context, in *endpoint.DescribeEndpointRequest) (*endpoint.Endpoint, error) {
	query := datasource.DBFromCtx(ctx)

	ins := &endpoint.Endpoint{}
	if err := query.First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("endpoint %d not found", in.Id)
		}
		return nil, err
	}

	return ins, nil
}
