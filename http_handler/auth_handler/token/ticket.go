package token

import (
	"github.com/getto-systems/project-example-id/basic"
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
	Parse(RenewToken, basic.Path) (basic.Ticket, error)
	RenewToken(basic.Ticket) (RenewToken, error)
	AppToken(basic.Ticket) (AppToken, error)
}
