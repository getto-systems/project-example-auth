package client

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	PasswordReset struct {
		Client
	}
	PasswordResetHandler interface {
		CreateSessionRequest() (request.Request, user.Login, error)
		CreateSessionResponse(password_reset.Session, error)

		SendTokenResponse(error)

		GetStatusRequest() (request.Request, user.Login, password_reset.Session, error)
		GetStatusResponse(password_reset.Destination, password_reset.Status, error)

		ResetRequest() (request.Request, user.Login, password_reset.Token, password.RawPassword, error)
		ResetResponse(error)
	}
)

func NewPasswordReset(client Client) PasswordReset {
	return PasswordReset{Client: client}
}

func (client PasswordReset) CreateSession(handler PasswordResetHandler) {
	session, err := client.createSession(handler)
	client.handleCredentialError(err)
	handler.CreateSessionResponse(session, err)
}
func (client PasswordReset) createSession(handler PasswordResetHandler) (_ password_reset.Session, err error) {
	request, login, err := handler.CreateSessionRequest()
	if err != nil {
		return
	}

	user, err := client.user.GetUser(request, login)
	if err != nil {
		return
	}

	session, dest, token, err := client.passwordReset.CreateSession(request, user, login)
	if err != nil {
		return
	}

	// job の追加は一番最後 : この後にエラーが発生した場合、再試行により job が 2重に登録されてしまう
	err = client.passwordReset.PushSendTokenJob(request, session, dest, token)
	if err != nil {
		return
	}

	return session, nil
}

func (client PasswordReset) SendToken(handler PasswordResetHandler) {
	err := client.sendToken(handler)
	client.handleCredentialError(err)
	handler.SendTokenResponse(err)
}
func (client PasswordReset) sendToken(handler PasswordResetHandler) error {
	return client.passwordReset.SendToken()
}

func (client PasswordReset) GetStatus(handler PasswordResetHandler) {
	dest, status, err := client.getStatus(handler)
	client.handleCredentialError(err)
	handler.GetStatusResponse(dest, status, err)
}
func (client PasswordReset) getStatus(handler PasswordResetHandler) (_ password_reset.Destination, _ password_reset.Status, err error) {
	request, login, session, err := handler.GetStatusRequest()
	if err != nil {
		return
	}

	return client.passwordReset.GetStatus(request, login, session)
}

func (client PasswordReset) Reset(handler PasswordResetHandler) {
	credential, err := client.reset(handler)
	client.handleCredential(credential, err)
	handler.ResetResponse(err)
}
func (client PasswordReset) reset(handler PasswordResetHandler) (_ credential.Credential, err error) {
	request, login, token, newPassword, err := handler.ResetRequest()
	if err != nil {
		return
	}

	user, session, exp, err := client.passwordReset.Validate(request, login, token)
	if err != nil {
		return
	}

	err = client.password.Change(request, user, newPassword)
	if err != nil {
		return
	}

	err = client.passwordReset.CloseSession(request, session)
	if err != nil {
		return
	}

	nonce, expires, err := client.ticket.Register(request, user, exp)
	if err != nil {
		return
	}

	return client.issueCredential(request, user, nonce, expires)
}
