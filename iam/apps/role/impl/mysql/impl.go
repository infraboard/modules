package mysql

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/role"
)

func init() {
	ioc.Controller().Registry(&RoleServiceImpl{})
}

type RoleServiceImpl struct {
	ioc.ObjectImpl
}

func (i *RoleServiceImpl) Init() error {
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&role.Role{}, &role.ApiPermission{}, &role.MenuPermission{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *RoleServiceImpl) Name() string {
	return role.AppName
}
