package credential

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"
)

type (
	Credential struct {
		ticket       Ticket
		apiToken     ApiToken
		contentToken ContentToken
		expires      expiration.Expires
	}
)

func NewCredential(ticket Ticket, apiToken ApiToken, contentToken ContentToken, expires expiration.Expires) Credential {
	return Credential{
		ticket:       ticket,
		apiToken:     apiToken,
		contentToken: contentToken,
		expires:      expires,
	}
}
func (credential Credential) Ticket() Ticket {
	return credential.ticket
}
func (credential Credential) ApiToken() ApiToken {
	return credential.apiToken
}
func (credential Credential) ContentToken() ContentToken {
	return credential.contentToken
}
func (credential Credential) Expires() expiration.Expires {
	return credential.expires
}
