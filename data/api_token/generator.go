package api_token

type (
	TicketNonceGenerator interface {
		GenerateNonce() (TicketNonce, error)
	}
)
