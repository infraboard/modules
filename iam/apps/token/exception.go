package token

import "github.com/infraboard/mcube/v2/exception"

func NewAuthFailed(format string, a ...any) *exception.ApiException {
	return exception.NewApiException(5000, "Auth Failed").WithMessagef(format, a...)
}

func NewPermissionDeny(format string, a ...any) *exception.ApiException {
	return exception.NewApiException(5100, "Permission Deny").WithMessagef(format, a...)
}

func NewAccessTokenExpired(format string, a ...any) *exception.ApiException {
	return exception.NewApiException(5001, "Access Token Expired").WithMessagef(format, a...)
}

func NewRefreshTokenExpired(format string, a ...any) *exception.ApiException {
	return exception.NewApiException(5002, "Refresh Token Expired").WithMessagef(format, a...)
}
