package nonce_generator

import (
	"github.com/google/uuid"

	infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data/credential"
)

type NonceGenerator struct {
}

func NewNonceGenerator() NonceGenerator {
	return NonceGenerator{}
}

func (gen NonceGenerator) gen() infra.TicketNonceGenerator {
	return gen
}

func (NonceGenerator) GenerateNonce() (_ credential.TicketNonce, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return credential.TicketNonce(id.String()), nil
}
