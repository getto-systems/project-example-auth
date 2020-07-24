package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type ResetLogger interface {
	TryToCreateResetSession(data.Request, Login, data.Expires)
	FailedToCreateResetSession(data.Request, Login, data.Expires, error)
	CreatedResetSession(data.Request, Login, data.Expires, data.User, ResetSession)

	TryToGetResetStatus(data.Request, ResetSession)
	FailedToGetResetStatus(data.Request, ResetSession, error)

	TryToValidateResetToken(data.Request)
	FailedToValidateResetToken(data.Request, error)
	AuthedByResetToken(data.Request, data.User)
}
