package data

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	Credential struct {
		ticket       ticket.Ticket
		apiToken     api_token.ApiToken
		contentToken api_token.ContentToken
		expires      time.Expires
	}
)

func NewCredential(ticket ticket.Ticket, apiToken api_token.ApiToken, contentToken api_token.ContentToken, expires time.Expires) Credential {
	return Credential{
		ticket:       ticket,
		apiToken:     apiToken,
		contentToken: contentToken,
		expires:      expires,
	}
}
func (credential Credential) Ticket() ticket.Ticket {
	return credential.ticket
}
func (credential Credential) ApiToken() api_token.ApiToken {
	return credential.apiToken
}
func (credential Credential) ContentToken() api_token.ContentToken {
	return credential.contentToken
}
func (credential Credential) Expires() time.Expires {
	return credential.expires
}
