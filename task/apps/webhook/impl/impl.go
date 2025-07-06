package impl

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/task/apps/webhook"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&WebHookServiceImpl{})
}

var _ webhook.Service = (*WebHookServiceImpl)(nil)

type WebHookServiceImpl struct {
	ioc.ObjectImpl

	log *zerolog.Logger
}

func (i *WebHookServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&webhook.WebHook{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *WebHookServiceImpl) Name() string {
	return webhook.APP_NAME
}
