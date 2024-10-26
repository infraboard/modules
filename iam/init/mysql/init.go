package mysql

import (
	_ "github.com/infraboard/modules/iam/apps/endpoint/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/namespace/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/policy/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/role/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/token/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/token/issuer"
	_ "github.com/infraboard/modules/iam/apps/user/impl/mysql"
	_ "github.com/infraboard/modules/iam/apps/view/impl/mysql"
	_ "github.com/infraboard/modules/iam/cmd"
)
