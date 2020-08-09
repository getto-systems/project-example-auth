package credential

import (
	"github.com/getto-systems/project-example-id/data/credential"
)

type (
	TicketNonceGenerator interface {
		GenerateNonce() (credential.TicketNonce, error)
	}
)
