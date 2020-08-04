package ticket

import (
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	Token  []byte
	Nonce  string
	Ticket struct {
		token Token
		nonce Nonce
	}

	ExpirationParam struct {
		Expires     time.Second
		ExtendLimit time.Second
	}
)

func NewTicket(token Token, nonce Nonce) Ticket {
	return Ticket{
		token: token,
		nonce: nonce,
	}
}
func (ticket Ticket) Token() Token {
	return ticket.token
}
func (ticket Ticket) Nonce() Nonce {
	return ticket.nonce
}