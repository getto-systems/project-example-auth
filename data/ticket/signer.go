package ticket

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketSigner interface {
		Sign(user.User, Nonce, time.Expires) (Token, error)
	}

	TicketParser interface {
		Parse(Token) (user.User, Nonce, error)
	}

	TicketSign interface {
		TicketSigner
		TicketParser
	}
)
