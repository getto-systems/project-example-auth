package ticket

type (
	NonceGenerator interface {
		GenerateNonce() (Nonce, error)
	}
)
