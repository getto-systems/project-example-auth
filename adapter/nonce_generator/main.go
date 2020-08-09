package nonce_generator

import (
	"github.com/google/uuid"

	"github.com/getto-systems/project-example-id/data/api_token"
)

type NonceGenerator struct {
}

func NewNonceGenerator() NonceGenerator {
	return NonceGenerator{}
}

func (gen NonceGenerator) gen() api_token.TicketNonceGenerator {
	return gen
}

func (NonceGenerator) GenerateNonce() (_ api_token.TicketNonce, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return api_token.TicketNonce(id.String()), nil
}
