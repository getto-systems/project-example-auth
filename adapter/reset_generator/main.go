package reset_generator

import (
	"github.com/google/uuid"

	"github.com/getto-systems/project-example-id/password"
)

type ResetGenerator struct {
}

func NewResetGenerator() ResetGenerator {
	return ResetGenerator{}
}

func (gen ResetGenerator) gen() password.ResetGenerator {
	return gen
}

func (ResetGenerator) Generate() (password.ResetID, password.ResetToken, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", "", err
	}

	token, err := uuid.NewUUID()
	if err != nil {
		return "", "", err
	}

	return password.ResetID(id.String()), password.ResetToken(token.String()), nil
}
