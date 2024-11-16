package mysql

import (
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/types"
	"github.com/infraboard/modules/system/apps/config"
	"gorm.io/gorm"
)

// 添加配置
func (i *ConfigServiceImpl) AddConfig(ctx context.Context, in *config.AddConfigRequest) (*types.Set[*config.ConfigItem], error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	set := types.New[*config.ConfigItem]()
	datasource.DBFromCtx(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range in.Items {
			ins := config.NewConfigItem()
			item := in.Items[i]
			ins.KVItem = *item
			if err := ins.Encrypt(); err != nil {
				return err
			}
			if err := tx.Save(ins).Error; err != nil {
				return err
			}
			set.Add(ins)
		}
		return nil
	})

	return set, nil
}

// 查询配置项
func (i *ConfigServiceImpl) QueryConfig(ctx context.Context, in *config.QueryConfigRequest) (*types.Set[*config.ConfigItem], error) {
	set := types.New[*config.ConfigItem]()

	query := datasource.DBFromCtx(ctx).Model(&config.ConfigItem{})
	err := query.
		Order("created_at desc").
		Find(&set.Items).
		Error
	if err != nil {
		return nil, err
	}
	return set, nil
}

// 查询配置详情
func (i *ConfigServiceImpl) DescribeConfig(ctx context.Context, in *config.DescribeConfigRequest) (*config.ConfigItem, error) {
	query := datasource.DBFromCtx(ctx)

	switch in.DescribeBy {
	case config.DESCRIBE_BY_ID:
		query = query.Where("id = ?", in.DescribeValue)
	case config.DESCRIBE_BY_KEY:
		query = query.Where("`key` = ?", in.DescribeValue)
	}

	ins := &config.ConfigItem{}
	if err := query.First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFound("config %s not found", in.DescribeValue)
		}
		return nil, err
	}

	if err := ins.Decrypt(); err != nil {
		return nil, err
	}

	return ins, nil
}

// 更新配置
func (i *ConfigServiceImpl) UpdateConfig(ctx context.Context, in *config.UpdateConfigRequest) (*config.ConfigItem, error) {
	descReq := config.NewDescribeConfigRequestById(in.Id)
	ins, err := i.DescribeConfig(ctx, descReq)
	if err != nil {
		return nil, err
	}
	ins.KVItem = in.KVItem

	if err := ins.Encrypt(); err != nil {
		return nil, err
	}

	return ins, datasource.DBFromCtx(ctx).Where("id = ?", in.Id).Updates(ins).Error
}
