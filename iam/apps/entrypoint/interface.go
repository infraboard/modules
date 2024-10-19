package entrypoint

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "entrypoint"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
}
