package mysql

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/entrypoint"
)

func init() {
	ioc.Controller().Registry(&EntryPointServiceImpl{})
}

// 他是user service 服务的控制器
type EntryPointServiceImpl struct {
	ioc.ObjectImpl
}

func (i *EntryPointServiceImpl) Init() error {
	// 自动创建表
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&entrypoint.Endpoint{})
		if err != nil {
			return err
		}
	}
	return nil
}

// 定义托管到Ioc里面的名称
func (i *EntryPointServiceImpl) Name() string {
	return entrypoint.AppName
}
