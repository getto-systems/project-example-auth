package infra

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/user"
)

type (
	TicketSigner interface {
		Sign(user.User, credential.TicketNonce, credential.Expires) (credential.TicketSignature, error)
	}
	TicketParser interface {
		Parse(credential.TicketSignature) (user.User, credential.TicketNonce, error)
	}
	TicketSign interface {
		TicketSigner
		TicketParser
	}

	ApiTokenSigner interface {
		Sign(user.User, credential.ApiRoles, credential.Expires) (credential.ApiToken, error)
	}

	ContentTokenSigner interface {
		Sign(credential.Expires) (credential.ContentToken, error)
	}
)
