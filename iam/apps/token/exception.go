package token

import "github.com/infraboard/mcube/v2/exception"

func NewAuthFailed(format string, a ...any) *exception.APIException {
	return exception.NewAPIException(5000, "Auth Failed").WithMessagef(format, a...)
}

func NewPermissionDeny(format string, a ...any) *exception.APIException {
	return exception.NewAPIException(5100, "Permission Deny").WithMessagef(format, a...)
}

func NewTokenExpired(format string, a ...any) *exception.APIException {
	return exception.NewAPIException(5001, "Token Expired").WithMessagef(format, a...)
}
