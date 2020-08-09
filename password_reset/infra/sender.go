package infra

import (
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	TokenSender interface {
		SendToken(password_reset.Destination, password_reset.Token) error
	}
)
