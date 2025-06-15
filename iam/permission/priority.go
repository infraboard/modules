package permission

import "github.com/infraboard/mcube/v2/ioc/config/gorestful"

func GetCheckerPriority() int {
	return gorestful.Priority() - 1
}
