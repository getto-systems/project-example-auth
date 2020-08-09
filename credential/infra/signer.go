package infra

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/user"
)

type (
	TicketSigner interface {
		Sign(user.User, credential.TicketNonce, expiration.Expires) (credential.TicketSignature, error)
	}
	TicketParser interface {
		Parse(credential.TicketSignature) (user.User, credential.TicketNonce, error)
	}
	TicketSign interface {
		TicketSigner
		TicketParser
	}

	ApiTokenSigner interface {
		Sign(user.User, credential.ApiRoles, expiration.Expires) (credential.ApiToken, error)
	}

	ContentTokenSigner interface {
		Sign(expiration.Expires) (credential.ContentToken, error)
	}
)
