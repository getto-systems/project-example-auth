package static_core

import (
	"github.com/getto-systems/project-example-auth/y_static"
)

func (action action) GetEnv() (_ static.Env, err error) {
	env, err := action.envs.FindEnv()
	if err != nil {
		return
	}

	return env, nil
}
