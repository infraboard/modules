package apps

import (
	_ "github.com/infraboard/modules/task/apps/cronjob/api"
	_ "github.com/infraboard/modules/task/apps/cronjob/impl"
	_ "github.com/infraboard/modules/task/apps/event/impl"
	_ "github.com/infraboard/modules/task/apps/task/api"
	_ "github.com/infraboard/modules/task/apps/task/impl"
	_ "github.com/infraboard/modules/task/apps/webhook/impl"
)
