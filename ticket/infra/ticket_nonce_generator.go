package infra

import (
	"github.com/getto-systems/project-example-id/credential"
)

type (
	TicketNonceGenerator interface {
		GenerateTicketNonce() (credential.TicketNonce, error)
	}
)
