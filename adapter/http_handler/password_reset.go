package http_handler

import (
	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
)

type PasswordReset struct {
	Handler
}

func NewPasswordReset(handler Handler) PasswordReset {
	return PasswordReset{Handler: handler}
}

func (handler PasswordReset) handler() client.PasswordResetHandler {
	return handler
}

func (handler PasswordReset) CreateSessionRequest() (_ request.Request, _ user.Login, err error) {
	type body struct {
		LoginID string `json:"login_id"`
	}

	var input body
	err = handler.parseBody(&input)
	if err != nil {
		return
	}

	login := user.NewLogin(user.LoginID(input.LoginID))

	return handler.newRequest(), login, nil
}
func (handler PasswordReset) CreateSessionResponse(session password_reset.Session, err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	type body struct {
		SessionID string `json:"session_id"`
	}

	handler.ok(body{
		SessionID: string(session.ID()),
	})
}

func (handler PasswordReset) SendTokenResponse(err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}

func (handler PasswordReset) GetStatusRequest() (_ request.Request, _ user.Login, _ password_reset.Session, err error) {
	type body struct {
		LoginID   string `json:"login_id"`
		SessionID string `json:"session_id"`
	}

	var input body
	err = handler.parseBody(&input)
	if err != nil {
		return
	}

	login := user.NewLogin(user.LoginID(input.LoginID))
	session := password_reset.NewSession(password_reset.SessionID(input.SessionID))

	return handler.newRequest(), login, session, nil
}
func (handler PasswordReset) GetStatusResponse(dest password_reset.Destination, status password_reset.Status, err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	// TODO dest をちゃんと返す
	// TODO status をちゃんと返す
	handler.ok("STATUS")
}

func (handler PasswordReset) ResetRequest() (_ request.Request, _ user.Login, _ password_reset.Token, _ password.RawPassword, err error) {
	type body struct {
		LoginID  string `json:"login_id"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var input body
	err = handler.parseBody(&input)
	if err != nil {
		return
	}

	login := user.NewLogin(user.LoginID(input.LoginID))
	token := password_reset.Token(input.Token)
	raw := password.RawPassword(input.Password)

	return handler.newRequest(), login, token, raw, nil
}
func (handler PasswordReset) ResetResponse(err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}
