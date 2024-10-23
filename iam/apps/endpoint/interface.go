package endpoint

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "endpoint"
)

func GetService() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
}
