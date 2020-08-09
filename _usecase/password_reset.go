package _usecase

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	PasswordReset struct {
		Usecase
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

func NewPasswordReset(u Usecase) PasswordReset {
	return PasswordReset{Usecase: u}
}

func (u PasswordReset) CreateSession(handler PasswordResetHandler) {
	session, err := u.createSession(handler)
	u.handleCredentialError(err)
	handler.CreateSessionResponse(session, err)
}
func (u PasswordReset) createSession(handler PasswordResetHandler) (_ password_reset.Session, err error) {
	request, login, err := handler.CreateSessionRequest()
	if err != nil {
		return
	}

	user, err := u.user.GetUser(request, login)
	if err != nil {
		return
	}

	session, dest, token, err := u.passwordReset.CreateSession(request, user, login)
	if err != nil {
		return
	}

	// job の追加は一番最後 : この後にエラーが発生した場合、再試行により job が 2重に登録されてしまう
	err = u.passwordReset.PushSendTokenJob(request, session, dest, token)
	if err != nil {
		return
	}

	return session, nil
}

func (u PasswordReset) SendToken(handler PasswordResetHandler) {
	err := u.sendToken(handler)
	u.handleCredentialError(err)
	handler.SendTokenResponse(err)
}
func (u PasswordReset) sendToken(handler PasswordResetHandler) error {
	return u.passwordReset.SendToken()
}

func (u PasswordReset) GetStatus(handler PasswordResetHandler) {
	dest, status, err := u.getStatus(handler)
	u.handleCredentialError(err)
	handler.GetStatusResponse(dest, status, err)
}
func (u PasswordReset) getStatus(handler PasswordResetHandler) (_ password_reset.Destination, _ password_reset.Status, err error) {
	request, login, session, err := handler.GetStatusRequest()
	if err != nil {
		return
	}

	return u.passwordReset.GetStatus(request, login, session)
}

func (u PasswordReset) Reset(handler PasswordResetHandler) {
	credential, err := u.reset(handler)
	u.handleCredential(credential, err)
	handler.ResetResponse(err)
}
func (u PasswordReset) reset(handler PasswordResetHandler) (_ credential.Credential, err error) {
	request, login, token, newPassword, err := handler.ResetRequest()
	if err != nil {
		return
	}

	user, session, exp, err := u.passwordReset.Validate(request, login, token)
	if err != nil {
		return
	}

	err = u.password.Change(request, user, newPassword)
	if err != nil {
		return
	}

	err = u.passwordReset.CloseSession(request, session)
	if err != nil {
		return
	}

	nonce, expires, err := u.ticket.Register(request, user, exp)
	if err != nil {
		return
	}

	return u.issueCredential(request, user, nonce, expires)
}
