package token

import (
	"github.com/getto-systems/project-example-id/user"
)

type TicketTokener interface {
	Parse(TicketToken, user.Path) (user.Ticket, error)
	Token(user.Ticket) (TicketToken, error)
	Info(user.Ticket) (TicketInfo, error)
}

type TicketToken []byte // use cookie value : e.g. JWT
type TicketInfo []byte  // use http response value : e.g. json data '{user_id, roles, token}'
