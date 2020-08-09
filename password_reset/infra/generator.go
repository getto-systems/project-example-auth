package infra

import (
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	SessionGenerator interface {
		GenerateSession() (password_reset.SessionID, password_reset.Token, error)
	}
)
