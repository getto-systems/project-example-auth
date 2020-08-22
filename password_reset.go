package auth

import (
	"errors"

	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

var (
	ErrInvalidPasswordReset = errors.New("invalid-password-reset")
	ErrClosedPasswordReset  = errors.New("closed-password-reset")
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
	handler.CreateSessionResponse(u.createSession(handler))
}
func (u PasswordReset) createSession(handler PasswordResetHandler) (_ password_reset.Session, err error) {
	request, login, err := handler.CreateSessionRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	resetUser, err := u.user.GetUser(request, login)
	if err != nil {
		switch err {
		case user.ErrGetUserNotFoundUser:
			err = ErrInvalidPasswordReset

		default:
			err = ErrServerError
		}
		return
	}

	session, dest, token, err := u.passwordReset.CreateSession(request, resetUser, login)
	if err != nil {
		switch err {
		case password_reset.ErrCreateSessionNotFoundDestination:
			err = ErrInvalidPasswordReset

		default:
			err = ErrServerError
		}
		return
	}

	// job の追加は一番最後 : この後にエラーが発生した場合、再試行により job が 2重に登録されてしまう
	err = u.passwordReset.PushSendTokenJob(request, session, dest, token)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	return session, nil
}

func (u PasswordReset) SendToken(handler PasswordResetHandler) {
	handler.SendTokenResponse(u.sendToken(handler))
}
func (u PasswordReset) sendToken(handler PasswordResetHandler) (err error) {
	err = u.passwordReset.SendToken()
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	return nil
}

func (u PasswordReset) GetStatus(handler PasswordResetHandler) {
	handler.GetStatusResponse(u.getStatus(handler))
}
func (u PasswordReset) getStatus(handler PasswordResetHandler) (_ password_reset.Destination, _ password_reset.Status, err error) {
	request, login, session, err := handler.GetStatusRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	dest, status, err := u.passwordReset.GetStatus(request, login, session)
	if err != nil {
		switch err {
		case password_reset.ErrGetStatusMatchFailedLogin,
			password_reset.ErrGetStatusNotFoundSession:

			err = ErrInvalidPasswordReset

		default:
			err = ErrServerError
		}
		return
	}

	return dest, status, nil
}

func (u PasswordReset) Reset(handler PasswordResetHandler) {
	handler.ResetResponse(u.reset(handler))
}
func (u PasswordReset) reset(handler PasswordResetHandler) (err error) {
	request, login, token, newPassword, err := handler.ResetRequest()
	if err != nil {
		switch err {
		default:
			err = ErrBadRequest
		}
		return
	}

	user, session, extendSecond, err := u.passwordReset.Validate(request, login, token)
	if err != nil {
		switch err {
		case password_reset.ErrValidateAlreadyClosed:
			err = ErrClosedPasswordReset

		case password_reset.ErrValidateAlreadyExpired,
			password_reset.ErrValidateMatchFailedLogin,
			password_reset.ErrValidateNotFoundSession:

			err = ErrInvalidPasswordReset

		default:
			err = ErrServerError
		}
		return
	}

	err = u.password.Change(request, user, newPassword)
	if err != nil {
		switch err {
		case password.ErrCheckLengthEmpty,
			password.ErrCheckLengthTooLong:

			err = ErrInvalidNewPassword

		default:
			err = ErrServerError
		}
		return
	}

	err = u.passwordReset.CloseSession(request, session)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	ticket, err := u.ticket.Register(request, user, extendSecond)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	err = u.issueCredential(request, ticket)
	if err != nil {
		return
	}

	return nil
}
