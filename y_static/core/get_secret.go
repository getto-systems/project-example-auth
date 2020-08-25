package static_core

import (
	"github.com/getto-systems/project-example-auth/y_static"
)

func (action action) GetSecret(name string) (_ static.Secret, err error) {
	secret, err := action.secrets.FindSecret(name)
	if err != nil {
		return
	}

	return secret, nil
}
