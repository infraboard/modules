package token

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/infraboard/mcube/v2/tools/pretty"
)

func GetAccessTokenFromHTTP(r *http.Request) string {
	// 先从Token中获取
	tk := r.Header.Get(ACCESS_TOKEN_HEADER_NAME)

	// 1. 获取Token
	if tk == "" {
		cookie, err := r.Cookie(ACCESS_TOKEN_COOKIE_NAME)
		if err != nil {
			return ""
		}
		tk, _ = url.QueryUnescape(cookie.Value)
	}
	return tk
}

func GetRefreshTokenFromHTTP(r *http.Request) string {
	// 先从Token中获取
	tk := r.Header.Get(REFRESH_TOKEN_HEADER_NAME)
	return tk
}

func NewToken() *Token {
	tk := &Token{
		// 生产一个UUID的字符串
		AccessToken:  MakeBearer(24),
		RefreshToken: MakeBearer(32),
		IssueAt:      time.Now(),
		Extras:       map[string]string{},
	}

	return tk
}

type Token struct {
	// 在添加数据需要村的定义
	Id uint64 `json:"id" gorm:"column:id;type:uint;primary_key;"`
	// 用户来源
	Source SOURCE `json:"source" gorm:"column:source;type:tinyint(1);index"`
	// 颁发器
	Issuer string `json:"issuer" gorm:"column:issuer;type:varchar(100);index"`
	// 该Token是颁发
	UserId uint64 `json:"user_id" gorm:"column:user_id;index"`
	// 用户名
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(100);not null;index"`
	// 是不是管理员
	IsAdmin bool `json:"is_admin" gorm:"column:is_admin;type:tinyint(1)"`
	// 颁发给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"access_token" gorm:"column:access_token;type:varchar(100);not null;uniqueIndex"`
	// 访问令牌过期时间
	AccessTokenExpiredAt *time.Time `json:"access_token_expired_at" gorm:"column:access_token_expired_at;type:timestamp;index"`
	// 刷新Token
	RefreshToken string `json:"refresh_token" gorm:"column:refresh_token;type:varchar(100);not null;uniqueIndex"`
	// 刷新Token过期时间
	RefreshTokenExpiredAt *time.Time `json:"refresh_token_expired_at" gorm:"column:refresh_token_expired_at;type:timestamp;index"`
	// 创建时间
	IssueAt time.Time `json:"issue_at" gorm:"column:issue_at;type:timestamp;default:current_timestamp;not null;index;"`
	// 更新时间
	RefreshAt *time.Time `json:"refresh_at" gorm:"column:refresh_at;type:timestamp;"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) IsAccessTokenExpired() error {
	if t.AccessTokenExpiredAt != nil {
		expiredSeconds := time.Since(*t.AccessTokenExpiredAt).Seconds()
		if expiredSeconds > 0 {
			return NewAccessTokenExpired("access token %s 过期了 %f秒",
				t.AccessToken, expiredSeconds)
		}
	}

	return nil
}

func (t *Token) IsRreshTokenExpired() error {
	if t.RefreshTokenExpiredAt != nil {
		expiredSeconds := time.Since(*t.RefreshTokenExpiredAt).Seconds()
		if expiredSeconds > 0 {
			return NewRefreshTokenExpired("refresh token %s 过期了 %f秒",
				t.RefreshToken, expiredSeconds)
		}
	}

	return nil
}

func (t *Token) SetExpiredAtByDuration(duration time.Duration, refreshMulti uint) {
	t.SetAccessTokenExpiredAt(time.Now().Add(duration))
	t.SetRefreshTokenExpiredAt(time.Now().Add(duration * time.Duration(refreshMulti)))
}

func (t *Token) SetAccessTokenExpiredAt(v time.Time) {
	t.AccessTokenExpiredAt = &v
}

func (t *Token) SetRefreshAt(v time.Time) {
	t.RefreshAt = &v
}

func (t *Token) AccessTokenExpiredTTL() int {
	if t.AccessTokenExpiredAt != nil {
		return int(t.AccessTokenExpiredAt.Sub(t.IssueAt).Seconds())
	}
	return 0
}

func (t *Token) SetRefreshTokenExpiredAt(v time.Time) {
	t.RefreshTokenExpiredAt = &v
}

func (t *Token) String() string {
	return pretty.ToJSON(t)
}

func (t *Token) UserIdString() string {
	return fmt.Sprintf("%d", t.UserId)
}
