package client

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	PasswordLogin struct {
		Client
	}
	PasswordLoginHandler interface {
		LoginRequest() (request.Request, user.Login, password.RawPassword, error)
		LoginResponse(error)
	}
)

func NewPasswordLogin(client Client) PasswordLogin {
	return PasswordLogin{Client: client}
}

func (client PasswordLogin) Login(handler PasswordLoginHandler) {
	credential, err := client.login(handler)
	client.handleCredential(credential, err)
	handler.LoginResponse(err)
}
func (client PasswordLogin) login(handler PasswordLoginHandler) (_ credential.Credential, err error) {
	request, login, raw, err := handler.LoginRequest()
	if err != nil {
		return
	}

	user, err := client.user.getUser.Get(request, login)
	if err != nil {
		return
	}

	exp, err := client.password.validate.Validate(request, user, raw)
	if err != nil {
		return
	}

	nonce, expires, err := client.ticket.register.Register(request, user, exp)
	if err != nil {
		return
	}

	return client.issueCredential(request, user, nonce, expires)
}
