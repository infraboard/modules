package impl

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/iam/apps/label"
)

func init() {
	ioc.Controller().Registry(&LabelServiceImpl{})
}

var _ label.Service = (*LabelServiceImpl)(nil)

type LabelServiceImpl struct {
	ioc.ObjectImpl
}

func (i *LabelServiceImpl) Init() error {
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&label.Label{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *LabelServiceImpl) Name() string {
	return label.APP_NAME
}
