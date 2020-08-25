package static_core

import (
	"github.com/getto-systems/project-example-auth/y_static/infra"

	"github.com/getto-systems/project-example-auth/y_static"
)

type (
	action struct {
		envs    infra.EnvRepository
		secrets infra.SecretRepository
	}
)

func NewAction(
	envs infra.EnvRepository,
	secrets infra.SecretRepository,
) static.Action {
	return action{
		envs:    envs,
		secrets: secrets,
	}
}
