package auth

import (
	"github.com/getto-systems/project-example-id/applog"

	"errors"
)

var (
	ErrUserPasswordNotFound    = errors.New("user password not found")
	ErrUserPasswordMatchFailed = errors.New("user password match failed")
	ErrUserAccessDenied        = errors.New("user access denied")
)

type Authenticator interface {
	Logger() applog.Logger
}
