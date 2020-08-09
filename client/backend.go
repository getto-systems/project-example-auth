package client

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	Backend struct {
		ticket        ticket.Action
		credential    credential.Action
		user          user.Action
		password      password.Action
		passwordReset password_reset.Action
	}
)

func NewBackend(
	ticket ticket.Action,
	credential credential.Action,
	user user.Action,
	password password.Action,
	passwordReset password_reset.Action,
) Backend {
	return Backend{
		ticket:        ticket,
		credential:    credential,
		user:          user,
		password:      password,
		passwordReset: passwordReset,
	}
}
