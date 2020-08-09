package infra

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketSigner interface {
		Sign(user.User, credential.TicketNonce, time.Expires) (credential.TicketSignature, error)
	}
	TicketParser interface {
		Parse(credential.TicketSignature) (user.User, credential.TicketNonce, error)
	}
	TicketSign interface {
		TicketSigner
		TicketParser
	}

	ApiTokenSigner interface {
		Sign(user.User, credential.ApiRoles, time.Expires) (credential.ApiToken, error)
	}

	ContentTokenSigner interface {
		Sign(time.Expires) (credential.ContentToken, error)
	}
)
