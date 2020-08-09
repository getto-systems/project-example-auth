package nonce_generator

import (
	"github.com/google/uuid"

	infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/credential"
)

type NonceGenerator struct {
}

func NewNonceGenerator() NonceGenerator {
	return NonceGenerator{}
}

func (gen NonceGenerator) gen() infra.TicketNonceGenerator {
	return gen
}

func (NonceGenerator) GenerateTicketNonce() (_ credential.TicketNonce, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return credential.TicketNonce(id.String()), nil
}
