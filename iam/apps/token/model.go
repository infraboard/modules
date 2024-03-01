package token

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/infraboard/modules/iam/apps/user"
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

func NewToken() *Token {
	return &Token{
		// 生产一个UUID的字符串
		AccessToken:           xid.New().String(),
		AccessTokenExpiredAt:  7200,
		RefreshToken:          xid.New().String(),
		RefreshTokenExpiredAt: 3600 * 24 * 7,
		CreatedAt:             time.Now().Unix(),
	}
}

type Token struct {
	// 该Token是颁发
	UserId string `json:"user_id"`
	// 人的名称， user_name
	UserName string `json:"username" gorm:"column:username"`
	// 办法给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"access_token"`
	// 过期时间(2h), 单位是秒
	AccessTokenExpiredAt int `json:"access_token_expired_at"`
	// 刷新Token
	RefreshToken string `json:"refresh_token"`
	// 刷新Token过期时间(7d)
	RefreshTokenExpiredAt int `json:"refresh_token_expired_at"`

	// 创建时间
	CreatedAt int64 `json:"created_at"`
	// 更新实现
	UpdatedAt int64 `json:"updated_at"`

	// 额外补充信息, gorm忽略处理
	Role user.Role `json:"role" gorm:"-"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) IsExpired() error {
	duration := time.Since(t.ExpiredTime())

	expiredSeconds := duration.Seconds()
	if expiredSeconds > 0 {
		return NewTokenExpired("token %s 过期了 %f秒",
			t.AccessToken, expiredSeconds)
	}

	return nil
}

// 计算Token的过期时间
func (t *Token) ExpiredTime() time.Time {
	return time.Unix(t.CreatedAt, 0).
		Add(time.Duration(t.AccessTokenExpiredAt) * time.Second)
}

func (u *Token) String() string {
	dj, _ := json.Marshal(u)
	return string(dj)
}
