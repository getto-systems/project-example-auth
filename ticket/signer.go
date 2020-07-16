package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type Signer interface {
	Verify(Ticket) (Nonce, data.User, data.Expires, error)
	Sign(Nonce, data.User, data.Expires) (Ticket, error)
}
