package token

import (
	"github.com/getto-systems/project-example-id/user"
)

type TicketSerializer interface {
	Parse(RenewToken, user.Path) (user.Ticket, error)
	RenewToken(user.Ticket) (RenewToken, error)
	AppToken(user.Ticket) (AppToken, error)
}

type RenewToken []byte
type AppToken struct {
	Token  string
	UserID user.UserID
	Roles  user.Roles
}
