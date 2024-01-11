package mysql

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/identity/apps/token"
	"github.com/infraboard/modules/identity/apps/user"
	"gorm.io/gorm"
)

func init() {
	ioc.Controller().Registry(&TokenServiceImpl{})
}

type TokenServiceImpl struct {
	ioc.ObjectImpl

	// db
	db *gorm.DB

	// 依赖User模块, 直接操作user模块的数据库(users)?
	// 这里需要依赖另一个业务领域: 用户管理领域
	user user.Service
}

func (i *TokenServiceImpl) Init() error {
	// db ioc
	i.db = datasource.DB().Debug()

	// 拿到的对象 在main 进行初始化好了
	i.user = ioc.Controller().Get(user.AppName).(user.Service)
	return nil
}

func (i *TokenServiceImpl) Name() string {
	return token.AppName
}
