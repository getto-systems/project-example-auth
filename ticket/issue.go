package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type IssueEventPublisher interface {
	IssueTicket(data.Request, data.User, data.Expires, data.ExtendLimit)
	IssueTicketFailed(data.Request, data.User, data.Expires, data.ExtendLimit, error)
}

type IssueDB interface {
	RegisterTransaction(Nonce, func(Nonce) error) (Nonce, error)
	RegisterTicket(Nonce, data.User, data.Expires, data.ExtendLimit) error
	NonceExists(Nonce) bool
}

type NonceGenerator interface {
	GenerateNonce() (Nonce, error)
}
