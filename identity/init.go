package identity

import (
	_ "github.com/infraboard/modules/identity/apps/token/api/gin"
	_ "github.com/infraboard/modules/identity/apps/token/impl/mysql"
	_ "github.com/infraboard/modules/identity/apps/user/api/gin"
	_ "github.com/infraboard/modules/identity/apps/user/impl/mysql"
)
