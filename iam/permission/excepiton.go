package permission

import "github.com/infraboard/mcube/v2/exception"

var (
	ErrUnauthorized = exception.NewUnauthorized("auth required")
)
