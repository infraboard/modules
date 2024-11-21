package mysql

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/config"
)

func init() {
	ioc.Controller().Registry(&ConfigServiceImpl{})
}

var _ config.Service = (*ConfigServiceImpl)(nil)

type ConfigServiceImpl struct {
	ioc.ObjectImpl
}

func (i *ConfigServiceImpl) Init() error {
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&config.ConfigItem{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ConfigServiceImpl) Name() string {
	return config.AppName
}
