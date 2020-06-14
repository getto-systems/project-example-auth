package auth

import (
	"github.com/getto-systems/project-example-id/applog"

	"errors"
)

var (
	ErrUserPasswordNotFound    = errors.New("user password not found")
	ErrUserPasswordDidNotMatch = errors.New("user password did not match")
	ErrUserAccessDenied        = errors.New("user access denied")
)

type Authenticator interface {
	Logger() applog.Logger
}
