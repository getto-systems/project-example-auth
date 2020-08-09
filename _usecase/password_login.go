package _usecase

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	PasswordLogin struct {
		Usecase
	}
	PasswordLoginHandler interface {
		LoginRequest() (request.Request, user.Login, password.RawPassword, error)
		LoginResponse(error)
	}
)

func NewPasswordLogin(u Usecase) PasswordLogin {
	return PasswordLogin{Usecase: u}
}

func (u PasswordLogin) Login(handler PasswordLoginHandler) {
	credential, err := u.login(handler)
	u.handleCredential(credential, err)
	handler.LoginResponse(err)
}
func (u PasswordLogin) login(handler PasswordLoginHandler) (_ credential.Credential, err error) {
	request, login, raw, err := handler.LoginRequest()
	if err != nil {
		return
	}

	user, err := u.user.GetUser(request, login)
	if err != nil {
		return
	}

	exp, err := u.password.Validate(request, user, raw)
	if err != nil {
		return
	}

	nonce, expires, err := u.ticket.Register(request, user, exp)
	if err != nil {
		return
	}

	return u.issueCredential(request, user, nonce, expires)
}
