package mysql

import (
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/modules/iam/apps/token"
	"github.com/infraboard/modules/iam/apps/user"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Controller().Registry(&TokenServiceImpl{
		AutoRefresh:     true,
		RereshTTLSecond: 1 * 60 * 60,
	})
}

var _ token.Service = (*TokenServiceImpl)(nil)

type TokenServiceImpl struct {
	ioc.ObjectImpl
	user user.Service
	log  *zerolog.Logger

	// 自动刷新, 直接刷新Token的过期时间，而不是生成一个新Token
	AutoRefresh bool `json:"auto_refresh" toml:"auto_refresh" yaml:"auto_refresh" env:"AUTO_REFRESH"`
	// 刷新TTL
	RereshTTLSecond uint64 `json:"refresh_ttl" toml:"refresh_ttl" yaml:"refresh_ttl" env:"REFRESH_TTL"`

	refreshDuration time.Duration
}

func (i *TokenServiceImpl) Init() error {
	i.log = log.Sub(i.Name())
	i.user = ioc.Controller().Get(user.AppName).(user.Service)
	i.refreshDuration = time.Duration(i.RereshTTLSecond) * time.Second

	if datasource.Get().AutoMigrate {
		err := datasource.DB().AutoMigrate(&token.Token{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *TokenServiceImpl) Name() string {
	return token.AppName
}
