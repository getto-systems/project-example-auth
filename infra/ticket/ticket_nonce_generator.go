package ticket

import (
	"github.com/getto-systems/project-example-id/data/credential"
)

type (
	TicketNonceGenerator interface {
		GenerateTicketNonce() (credential.TicketNonce, error)
	}
)
