package credential

type (
	TicketSignature []byte
	TicketNonce     string
	Ticket          struct {
		signature TicketSignature
		nonce     TicketNonce
	}

	ApiRoles     []string
	ApiSignature []byte
	ApiToken     struct {
		roles     ApiRoles
		signature ApiSignature
	}

	ContentKeyID     string
	ContentPolicy    string
	ContentSignature string
	ContentToken     struct {
		keyID     ContentKeyID
		policy    ContentPolicy
		signature ContentSignature
	}
)

func NewTicket(signature TicketSignature, nonce TicketNonce) Ticket {
	return Ticket{
		signature: signature,
		nonce:     nonce,
	}
}
func (ticket Ticket) Signature() TicketSignature {
	return ticket.signature
}
func (ticket Ticket) Nonce() TicketNonce {
	return ticket.nonce
}

func EmptyApiRoles() ApiRoles {
	return []string{}
}

func NewApiToken(roles ApiRoles, signature ApiSignature) ApiToken {
	return ApiToken{
		roles:     roles,
		signature: signature,
	}
}
func (token ApiToken) ApiRoles() ApiRoles {
	return token.roles
}
func (token ApiToken) Signature() ApiSignature {
	return token.signature
}

func NewContentToken(keyID ContentKeyID, policy ContentPolicy, signature ContentSignature) ContentToken {
	return ContentToken{
		keyID:     keyID,
		policy:    policy,
		signature: signature,
	}
}
func (token ContentToken) KeyID() ContentKeyID {
	return token.keyID
}
func (token ContentToken) Policy() ContentPolicy {
	return token.policy
}
func (token ContentToken) Signature() ContentSignature {
	return token.signature
}
