package password

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrResetUserNotFound  = errors.New("reset user not found")
	ErrResetTokenNotFound = errors.New("reset token not found")
)

type ResetLogger interface {
	TryToIssueReset(data.Request, Login, data.Expires)
	FailedToIssueReset(data.Request, Login, data.Expires, error)
	IssuedReset(data.Request, Login, data.Expires, Reset, data.User)

	TryToGetResetStatus(data.Request, Reset)
	FailedToGetResetStatus(data.Request, Reset, error)

	TryToValidateResetToken(data.Request)
	FailedToValidateResetToken(data.Request, error)
	AuthedByResetToken(data.Request, data.User)
}

type ResetDB interface {
	FilterUserByLogin(Login) ([]data.User, error)
	RegisterReset(ResetGenerator, data.User, data.RequestedAt, data.Expires) (Reset, ResetToken, error)

	FilterResetStatus(Reset) ([]ResetStatus, error)

	FilterResetUser(ResetToken) ([]ResetUser, error)
}

type ResetGenerator interface {
	Generate() (ResetID, ResetToken, error)
}
