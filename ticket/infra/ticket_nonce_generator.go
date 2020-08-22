package infra

import (
	"github.com/getto-systems/project-example-auth/credential"
)

type (
	TicketNonceGenerator interface {
		GenerateTicketNonce() (credential.TicketNonce, error)
	}
)
