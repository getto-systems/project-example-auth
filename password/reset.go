package password

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrResetUserNotFound  = errors.New("reset user not found")
	ErrResetTokenNotFound = errors.New("reset token not found")
)

type ResetEventPublisher interface {
	IssueReset(data.Request, Login, data.Expires)
	IssueResetFailed(data.Request, Login, data.Expires, error)
	IssuedReset(data.Request, Login, data.Expires, Reset, ResetToken)

	GetResetStatus(data.Request, Reset)
	GetResetStatusFailed(data.Request, Reset, error)

	ValidateResetToken(data.Request)
	ValidateResetTokenFailed(data.Request, error)
	AuthenticatedByResetToken(data.Request, data.User)
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
