package token

import "github.com/infraboard/mcube/v2/exception"

const (
	ACCESS_TOKEN_HEADER_NAME          = "Authorization"
	ACCESS_TOKEN_COOKIE_NAME          = "access_token"
	ACCESS_TOKEN_RESPONSE_HEADER_NAME = "X-OAUTH-TOKEN"
	REFRESH_TOKEN_HEADER_NAME         = "X-REFRUSH-TOKEN"
)

var (
	CTX_TOKEN_KEY = struct{}{}
)

var (
	CookieNotFound = exception.NewUnauthorized("cookie %s not found", ACCESS_TOKEN_COOKIE_NAME)
)

type SOURCE int

const (
	// 未知
	SOURCE_UNKNOWN SOURCE = iota
	// Web
	SOURCE_WEB
	// IOS
	SOURCE_IOS
	// ANDROID
	SOURCE_ANDROID
	// PC
	SOURCE_PC
	// API 调用
	SOURCE_API SOURCE = 10
)

type LOCK_TYPE int

const (
	// 用户退出登录
	LOCK_TYPE_REVOLK LOCK_TYPE = iota
	// 刷新Token过期, 回话中断
	LOCK_TYPE_TOKEN_EXPIRED
	// 异地登陆
	LOCK_TYPE_OTHER_PLACE_LOGGED_IN
	// 异常Ip登陆
	LOCK_TYPE_OTHER_IP_LOGGED_IN
)

type DESCRIBE_BY int

const (
	DESCRIBE_BY_ACCESS_TOKEN DESCRIBE_BY = iota
)
