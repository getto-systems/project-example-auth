package reset_session_generator

import (
	"github.com/google/uuid"

	"github.com/getto-systems/project-example-id/password"
)

type Generator struct {
}

func NewGenerator() Generator {
	return Generator{}
}

func (gen Generator) gen() password.ResetSessionGenerator {
	return gen
}

func (Generator) GenerateSession() (_ password.ResetSessionID, _ password.ResetToken, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	token, err := uuid.NewUUID()
	if err != nil {
		return
	}

	return password.ResetSessionID(id.String()), password.ResetToken(token.String()), nil
}
