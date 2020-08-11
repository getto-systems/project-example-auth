package credential

type (
	Credential struct {
		ticket       TicketToken
		apiToken     ApiToken
		contentToken ContentToken
	}
)

func NewCredential(ticket TicketToken, apiToken ApiToken, contentToken ContentToken) Credential {
	return Credential{
		ticket:       ticket,
		apiToken:     apiToken,
		contentToken: contentToken,
	}
}
func (credential Credential) TicketToken() TicketToken {
	return credential.ticket
}
func (credential Credential) ApiToken() ApiToken {
	return credential.apiToken
}
func (credential Credential) ContentToken() ContentToken {
	return credential.contentToken
}
