package infra

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	TicketSigner interface {
		Sign(user.User, credential.TicketNonce, credential.TicketExpires) (credential.TicketSignature, error)
	}
	TicketParser interface {
		Parse(credential.TicketSignature) (user.User, credential.TicketNonce, error)
	}
	TicketSign interface {
		TicketSigner
		TicketParser
	}

	ApiTokenSigner interface {
		Sign(user.User, credential.ApiRoles, credential.TokenExpires) (credential.ApiSignature, error)
	}

	ContentTokenSigner interface {
		Sign(credential.TokenExpires) (credential.ContentKeyID, credential.ContentPolicy, credential.ContentSignature, error)
	}
)
