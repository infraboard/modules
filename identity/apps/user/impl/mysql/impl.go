package imp

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/modules/identity/apps/user"
	"gorm.io/gorm"
)

func init() {
	ioc.Controller().Registry(&UserServiceImpl{})
}

// 显示声明接口实现的语言 都可以 明确约束接口的实现
var _ user.Service = (*UserServiceImpl)(nil)

func (i *UserServiceImpl) Init() error {
	i.db = datasource.DB().Debug()
	return nil
}

// 定义托管到Ioc里面的名称
func (i *UserServiceImpl) Name() string {
	return user.AppName
}

// 他是user service 服务的控制器
type UserServiceImpl struct {
	db *gorm.DB
	ioc.ObjectImpl
}
