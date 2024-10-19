package token

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
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

func NewToken(exipredDuration time.Duration) *Token {
	tk := &Token{
		// 生产一个UUID的字符串
		AccessToken:  xid.New().String(),
		RefreshToken: xid.New().String(),
		CreatedAt:    time.Now(),
		Extras:       map[string]string{},
	}
	if exipredDuration != 0 {
		tk.SetAccessTokenExpiredAt(time.Now().Add(exipredDuration))
		tk.SetRefreshTokenExpiredAt(time.Now().Add(exipredDuration * 4))
	}

	return tk
}

type Token struct {
	// 用户来源
	Source SOURCE `json:"source" gorm:"column:source;type:tinyint(1);index"`
	// 颁发器
	Issuer string `json:"issuer" gorm:"column:issuer;type:varchar(100);index"`
	// 该Token是颁发
	UserId uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;index"`
	// 用户名
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(100);not null;index"`
	// 是不是管理员
	IsAdmin bool `json:"is_admin" gorm:"column:is_admin;type:tinyint(1)"`
	// 办法给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"access_token" gorm:"column:access_token;type:varchar(100);not null;index"`
	// 访问令牌过期时间
	AccessTokenExpiredAt *time.Time `json:"access_token_expired_at" gorm:"column:access_token_expired_at;type:timestamp;index"`
	// 刷新Token
	RefreshToken string `json:"refresh_token" gorm:"column:refresh_token;type:varchar(100);not null;index"`
	// 刷新Token过期时间
	RefreshTokenExpiredAt *time.Time `json:"refresh_token_expired_at" gorm:"column:refresh_token_expired_at;type:timestamp;index"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp;not null;index;"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) IsExpired() error {
	if t.AccessTokenExpiredAt != nil {
		expiredSeconds := time.Since(*t.AccessTokenExpiredAt).Seconds()
		if expiredSeconds > 0 {
			return NewTokenExpired("token %s 过期了 %f秒",
				t.AccessToken, expiredSeconds)
		}
	}

	return nil
}

func (t *Token) SetAccessTokenExpiredAt(v time.Time) {
	t.AccessTokenExpiredAt = &v
}

func (t *Token) AccessTokenExpiredTTL() int {
	if t.AccessTokenExpiredAt != nil {
		return int(t.AccessTokenExpiredAt.Sub(t.CreatedAt).Seconds())
	}
	return 0
}

func (t *Token) SetRefreshTokenExpiredAt(v time.Time) {
	t.RefreshTokenExpiredAt = &v
}

func (u *Token) String() string {
	dj, _ := json.Marshal(u)
	return string(dj)
}
