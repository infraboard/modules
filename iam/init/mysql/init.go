package mysql

import (
	_ "github.com/infraboard/modules/iam/apps/token/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/token/issuer"
	_ "github.com/infraboard/modules/iam/apps/user/impl/mysql"
	_ "github.com/infraboard/modules/iam/cmd"
)
