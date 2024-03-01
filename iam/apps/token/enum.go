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
