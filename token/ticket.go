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
	Parse(RenewToken, basic.Path) (basic.TicketData, error)
	RenewToken(basic.TicketData) (RenewToken, error)
	AppToken(basic.TicketData) (AppToken, error)
}
