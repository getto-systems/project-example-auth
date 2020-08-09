package credential

type (
	Credential struct {
		ticket       Ticket
		apiToken     ApiToken
		contentToken ContentToken
		expires      Expires
	}
)

func NewCredential(ticket Ticket, apiToken ApiToken, contentToken ContentToken, expires Expires) Credential {
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
func (credential Credential) Expires() Expires {
	return credential.expires
}
