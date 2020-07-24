package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type resetter struct {
	logger    password.ResetLogger
	passwords password.PasswordRepository
	sessions  password.ResetSessionRepository
	exp       password.ResetSessionExpiration
	gen       password.ResetSessionGenerator
}

func newResetter(
	logger password.ResetLogger,
	passwords password.PasswordRepository,
	sessions password.ResetSessionRepository,
	exp password.ResetSessionExpiration,
	gen password.ResetSessionGenerator,
) resetter {
	return resetter{
		logger:    logger,
		passwords: passwords,
		sessions:  sessions,
		exp:       exp,
		gen:       gen,
	}
}

func (resetter resetter) createResetSession(request data.Request, login password.Login) (_ password.ResetSession, err error) {
	requestedAt := request.RequestedAt()
	expires := resetter.exp.Expires(requestedAt)

	resetter.logger.TryToCreateResetSession(request, login, expires)
	defer func() {
		if err != nil {
			resetter.logger.FailedToCreateResetSession(request, login, expires, err)
		}
	}()

	user, err := resetter.passwords.FindUser(login)
	if err != nil {
		return
	}

	data := password.NewResetSessionData(user, login, requestedAt, expires)

	// TODO token を notifier に渡す
	session, _, err := resetter.sessions.RegisterResetSession(resetter.gen, data)
	if err != nil {
		return
	}

	//resetter.notifier.SendResetToken(data, token)

	resetter.logger.CreatedResetSession(request, login, expires, user, session)
	return session, nil
}

func (resetter resetter) getResetStatus(request data.Request, session password.ResetSession) (_ password.ResetStatus, err error) {
	resetter.logger.TryToGetResetStatus(request, session)
	defer func() {
		if err != nil {
			resetter.logger.FailedToGetResetStatus(request, session, err)
		}
	}()

	return resetter.sessions.FindResetStatus(session)
}

func (resetter resetter) validate(request data.Request, login password.Login, token password.ResetToken) (_ data.User, err error) {
	resetter.logger.TryToValidateResetToken(request)
	defer func() {
		if err != nil {
			resetter.logger.FailedToValidateResetToken(request, err)
		}
	}()

	data, err := resetter.sessions.FindResetSession(token)
	if err != nil {
		return
	}

	if data.Login().ID() != login.ID() {
		err = password.ErrResetSessionLoginNotMatched
		return
	}

	if request.RequestedAt().Expired(data.Expires()) {
		err = password.ErrResetSessionAlreadyExpired
		return
	}

	resetter.logger.AuthedByResetToken(request, data.User())
	return data.User(), nil
}
