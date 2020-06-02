package token

import (
	"github.com/getto-systems/project-example-id/basic"
	"github.com/getto-systems/project-example-id/user"
)

type (
	RenewToken []byte
	AppToken   struct {
		Token  string
		UserID basic.UserID
		Roles  basic.Roles
	}
)

type TicketSerializer interface {
	Parse(RenewToken, basic.Path) (user.Ticket, error)
	RenewToken(user.Ticket) (RenewToken, error)
	AppToken(user.Ticket) (AppToken, error)
}
