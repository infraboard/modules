package token

const (
	ACCESS_TOKEN_HEADER_NAME          = "Authorization"
	ACCESS_TOKEN_COOKIE_NAME          = "access_token"
	ACCESS_TOKEN_GIN_KEY_NAME         = "access_token"
	ACCESS_TOKEN_RESPONSE_HEADER_NAME = "X-OAUTH-TOKEN"
	REFRESH_TOKEN_HEADER_NAME         = "X-REFRUSH-TOKEN"
)

var (
	CookieNotFound = NewAuthFailed("cookie %s not found", ACCESS_TOKEN_COOKIE_NAME)
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
