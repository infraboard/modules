package token

const (
	TOKEN_COOKIE_NAME  = "access_token"
	TOKEN_GIN_KEY_NAME = "access_token"
)

var (
	CookieNotFound = NewAuthFailed("cookie %s not found", TOKEN_COOKIE_NAME)
)
