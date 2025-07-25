package token

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/modules/iam/apps/policy"
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
	} else {
		// 处理 带格式: Bearer <Your API key>
		ft := strings.Split(tk, " ")
		if len(ft) > 1 {
			tk = ft[1]
		}
	}
	return tk
}

func GetTokenFromCtx(ctx context.Context) *Token {
	if v := ctx.Value(CTX_TOKEN_KEY); v != nil {
		return v.(*Token)
	}
	return nil
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
		Status:       NewStatus(),
		Extras:       map[string]string{},
	}

	return tk
}

type Token struct {
	// 在添加数据需要村的定义
	Id uint64 `json:"id" gorm:"column:id;type:uint;primary_key;"`
	// 用户来源
	Source SOURCE `json:"source" gorm:"column:source;type:tinyint(1);index" description:"用户来源"`
	// 颁发器
	Issuer string `json:"issuer" gorm:"column:issuer;type:varchar(100);index" description:"颁发器"`
	// 该Token属于哪个用户
	UserId uint64 `json:"user_id" gorm:"column:user_id;index" description:"持有该Token的用户Id"`
	// 用户名
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(100);not null;index" description:"持有该Token的用户名称"`
	// 是不是管理员
	IsAdmin bool `json:"is_admin" gorm:"column:is_admin;type:tinyint(1)" description:"是不是管理员"`

	// 令牌生效空间Ids
	policy.ResourceScope
	// 令牌生效空间名称
	NamespaceName string `json:"namespace_name" gorm:"column:namespace_name;type:varchar(100);index" description:"令牌所属空间"`
	// 颁发给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"access_token" gorm:"column:access_token;type:varchar(100);not null;uniqueIndex" description:"访问令牌"`
	// 访问令牌过期时间
	AccessTokenExpiredAt *time.Time `json:"access_token_expired_at" gorm:"column:access_token_expired_at;type:timestamp;index" description:"访问令牌的过期时间"`
	// 刷新Token
	RefreshToken string `json:"refresh_token" gorm:"column:refresh_token;type:varchar(100);not null;uniqueIndex" description:"刷新令牌"`
	// 刷新Token过期时间
	RefreshTokenExpiredAt *time.Time `json:"refresh_token_expired_at" gorm:"column:refresh_token_expired_at;type:timestamp;index" description:"刷新令牌的过期时间"`
	// 创建时间
	IssueAt time.Time `json:"issue_at" gorm:"column:issue_at;type:timestamp;default:current_timestamp;not null;index" description:"令牌颁发时间"`
	// 更新时间
	RefreshAt *time.Time `json:"refresh_at" gorm:"column:refresh_at;type:timestamp" description:"令牌刷新时间"`
	// 令牌状态
	Status *Status `json:"status" gorm:"embedded" modelDescription:"令牌状态"`
	// 其他扩展信息
	Extras map[string]string `json:"extras" gorm:"column:extras;serializer:json;type:json" description:"其他扩展信息"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) IsAccessTokenExpired() error {
	if t.AccessTokenExpiredAt != nil {
		expiredSeconds := time.Since(*t.AccessTokenExpiredAt).Seconds()
		if expiredSeconds > 0 {
			return exception.NewAccessTokenExpired("access token %s 过期了 %f秒",
				t.AccessToken, expiredSeconds)
		}
	}

	return nil
}

func (t *Token) IsRreshTokenExpired() error {
	if t.RefreshTokenExpiredAt != nil {
		expiredSeconds := time.Since(*t.RefreshTokenExpiredAt).Seconds()
		if expiredSeconds > 0 {
			return exception.NewRefreshTokenExpired("refresh token %s 过期了 %f秒",
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

func (t *Token) SetIssuer(issuer string) *Token {
	t.Issuer = issuer
	return t
}

func (t *Token) SetSource(source SOURCE) *Token {
	t.Source = source
	return t
}

func (t *Token) UserIdString() string {
	return fmt.Sprintf("%d", t.UserId)
}

func (t *Token) CheckRefreshToken(refreshToken string) error {
	if t.RefreshToken != refreshToken {
		return exception.NewPermissionDeny("refresh token not conrect")
	}
	return nil
}

func (t *Token) Lock(l LOCK_TYPE, reason string) {
	if t.Status == nil {
		t.Status = NewStatus()
	}
	t.Status.LockType = l
	t.Status.LockReason = reason
	t.Status.SetLockAt(time.Now())
}

func NewStatus() *Status {
	return &Status{}
}

type Status struct {
	// 冻结时间
	LockAt *time.Time `json:"lock_at" bson:"lock_at" gorm:"column:lock_at;type:timestamp;index" description:"冻结时间"`
	// 冻结类型
	LockType LOCK_TYPE `json:"lock_type" bson:"lock_type" gorm:"column:lock_type;type:tinyint(1)" description:"冻结类型 0:用户退出登录, 1:刷新Token过期, 回话中断, 2:异地登陆, 异常Ip登陆" enum:"0|1|2|3"`
	// 冻结原因
	LockReason string `json:"lock_reason" bson:"lock_reason" gorm:"column:lock_reason;type:text" description:"冻结原因"`
}

func (s *Status) SetLockAt(v time.Time) {
	s.LockAt = &v
}

func (s *Status) ToMap() map[string]any {
	return map[string]any{
		"lock_at":     s.LockAt,
		"lock_type":   s.LockType,
		"lock_reason": s.LockReason,
	}
}
