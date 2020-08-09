package credential

type (
	TicketNonceGenerator interface {
		GenerateNonce() (TicketNonce, error)
	}
)
