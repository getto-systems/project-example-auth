package core

import (
	"errors"
	"fmt"

	"github.com/getto-systems/project-example-id/password"
)

var (
	errPasswordEmpty   = newError("Password/Empty")
	errPasswordTooLong = newError("Password/TooLong")
)

var (
	errPasswordNotMatched       = newError("Password/NotMatched")
	errPasswordNotFoundLogin    = newError("Password/NotFound/Login")
	errPasswordNotFoundPassword = newError("Password/NotFound/Password")
)

var (
	errResetSessionNotFoundUser         = newError("ResetSession/NotFound/User")
	errResetSessionNotFoundResetStatus  = newError("ResetSession/NotFound/ResetStatus")
	errResetSessionNotFoundResetSession = newError("ResetSession/NotFound/ResetSession")
	errResetSessionLoginNotMatched      = newError("ResetSession/LoginNotMatched")
	errResetSessionAlreadyExpired       = newError("ResetSession/AlreadyExpired")
)

func newError(message string) error {
	return errors.New(fmt.Sprintf("Password/%s", message))
}

type (
	usecaseError struct{}
)

func newUsecaseError() (_ password.UsecaseError) {
	return
}

func (usecaseError) InputError(err error) bool {
	switch err {
	case
		errPasswordEmpty,
		errPasswordTooLong:
		return true
	default:
		return false
	}
}

func (usecaseError) AuthError(err error) bool {
	switch err {
	case
		errPasswordNotMatched,
		errPasswordNotFoundLogin,
		errPasswordNotFoundPassword:
		return true
	default:
		return false
	}
}

func (usecaseError) ResetError(err error) bool {
	switch err {
	case
		errResetSessionNotFoundUser,
		errResetSessionNotFoundResetStatus,
		errResetSessionNotFoundResetSession,
		errResetSessionLoginNotMatched,
		errResetSessionAlreadyExpired:
		return true
	default:
		return false
	}
}
