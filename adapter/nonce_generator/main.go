package nonce_generator

import (
	"github.com/google/uuid"

	"github.com/getto-systems/project-example-id/ticket"
)

type NonceGenerator struct {
}

func NewNonceGenerator() NonceGenerator {
	return NonceGenerator{}
}

func (gen NonceGenerator) gen() ticket.NonceGenerator {
	return gen
}

func (NonceGenerator) GenerateNonce() (_ ticket.Nonce, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return ticket.Nonce(id.String()), nil
}
