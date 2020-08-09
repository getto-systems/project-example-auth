package credential

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketSigner interface {
		Sign(user.User, TicketNonce, time.Expires) (TicketSignature, error)
	}
	TicketParser interface {
		Parse(TicketSignature) (user.User, TicketNonce, error)
	}
	TicketSign interface {
		TicketSigner
		TicketParser
	}

	ApiTokenSigner interface {
		Sign(user.User, ApiRoles, time.Expires) (ApiToken, error)
	}

	ContentTokenSigner interface {
		Sign(time.Expires) (ContentToken, error)
	}
)
