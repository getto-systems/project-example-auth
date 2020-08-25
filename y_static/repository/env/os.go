package static_repository_env

import (
	"os"

	"github.com/getto-systems/project-example-auth/y_static/infra"

	"github.com/getto-systems/project-example-auth/y_static"
)

type (
	OSEnv struct{}
)

func NewOSEnv() OSEnv {
	return OSEnv{}
}

func (env OSEnv) repo() infra.EnvRepository {
	return env
}

func (OSEnv) FindEnv() (_ static.Env, err error) {
	return static.Env{
		LogLevel:   os.Getenv("LOG_LEVEL"),
		SecretName: os.Getenv("GOOGLE_SECRET_MANAGER_SECRET_NAME"),
	}, nil
}
