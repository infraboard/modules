package iam

import (
	_ "github.com/infraboard/modules/iam/apps/token/api/gin"
	_ "github.com/infraboard/modules/iam/apps/token/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/user/api/gin"
	_ "github.com/infraboard/modules/iam/apps/user/impl/mysql"
)
