package infra

import (
	"github.com/getto-systems/project-example-auth/y_static"
)

type (
	EnvRepository interface {
		FindEnv() (static.Env, error)
	}

	SecretRepository interface {
		FindSecret(string) (static.Secret, error)
	}
)
