package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type registerer struct {
	logger    password.RegisterLogger
	passwords password.PasswordRepository
	gen       password.PasswordGenerator
}

func newRegisterer(
	logger password.RegisterLogger,
	passwords password.PasswordRepository,
	gen password.PasswordGenerator,
) registerer {
	return registerer{
		logger:    logger,
		gen:       gen,
		passwords: passwords,
	}
}

func (registerer registerer) getLogin(request data.Request, user data.User) (_ password.Login, err error) {
	registerer.logger.TryToGetLogin(request, user)
	defer func() {
		if err != nil {
			registerer.logger.FailedToGetLogin(request, user, err)
		}
	}()

	login, found, err := registerer.passwords.FindLogin(user)
	if err != nil {
		return
	}
	if !found {
		err = password.ErrPasswordNotFoundLogin
		return
	}

	return login, nil
}

func (registerer registerer) register(request data.Request, user data.User, raw password.RawPassword) (err error) {
	registerer.logger.TryToRegister(request, user)
	defer func() {
		if err != nil {
			registerer.logger.FailedToRegister(request, user, err)
		}
	}()

	err = raw.Check()
	if err != nil {
		return
	}

	hashed, err := registerer.gen.GeneratePassword(raw)
	if err != nil {
		return
	}

	err = registerer.passwords.RegisterPassword(user, hashed)
	if err != nil {
		return
	}

	registerer.logger.Registered(request, user)
	return nil
}
