package reset_session_generator

import (
	"github.com/google/uuid"

	infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
)

type Generator struct {
}

func NewGenerator() Generator {
	return Generator{}
}

func (gen Generator) gen() infra.SessionGenerator {
	return gen
}

func (Generator) GenerateSession() (_ password_reset.SessionID, _ password_reset.Token, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	token, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return password_reset.SessionID(id.String()), password_reset.Token(token.String()), nil
}
