package mysql

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/view"
)

func init() {
	ioc.Controller().Registry(&MenuServiceImpl{})
}

type MenuServiceImpl struct {
	ioc.ObjectImpl
}

func (i *MenuServiceImpl) Init() error {
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&view.Menu{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *MenuServiceImpl) Name() string {
	return view.AppName
}
