package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/iam/apps/config"
)

// 添加配置
func (i *ConfigServiceImpl) AddConfig(ctx context.Context, in *config.AddConfigRequest) (*config.ConfigItem, error) {
	return nil, nil
}

// 查询配置项
func (i *ConfigServiceImpl) QueryConfig(ctx context.Context, in *config.QueryConfigRequest) (*types.Set[*config.ConfigItem], error) {
	return nil, nil
}

// 查询配置详情
func (i *ConfigServiceImpl) DescribeConfig(ctx context.Context, in *config.DescribeConfigRequest) (*config.ConfigItem, error) {
	return nil, nil
}

// 更新配置
func (i *ConfigServiceImpl) UpdateConfig(ctx context.Context, in *config.UpdateConfigRequest) (*config.ConfigItem, error) {
	return nil, nil
}
