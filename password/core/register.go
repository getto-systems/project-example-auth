package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type registerer struct {
	logger password.RegisterLogger
	gen    password.Generator
	repo   registerRepository
}

func newRegisterer(
	logger password.RegisterLogger,
	db password.RegisterDB,
	gen password.Generator,
) registerer {
	return registerer{
		logger: logger,
		gen:    gen,
		repo:   newRegisterRepository(db, gen),
	}
}

func (registerer registerer) getLogin(request data.Request, user data.User) (login password.Login, err error) {
	registerer.logger.TryToGetLogin(request, user)
	defer func() {
		if err != nil {
			registerer.logger.FailedToGetLogin(request, user, err)
		}
	}()

	login, err = registerer.repo.findLogin(user)
	if err != nil {
		return login, err
	}

	return login, nil
}

func (registerer registerer) register(request data.Request, user data.User, password password.RawPassword) (err error) {
	registerer.logger.TryToRegister(request, user)
	defer func() {
		if err == nil {
			registerer.logger.Registered(request, user)
		} else {
			registerer.logger.FailedToRegister(request, user, err)
		}
	}()

	err = checkPassword(password)
	if err != nil {
		return err
	}

	hashed, err := registerer.gen.GeneratePassword(password)
	if err != nil {
		return err
	}

	return registerer.repo.register(user, hashed)
}

type registerRepository struct {
	db password.RegisterDB
}

func newRegisterRepository(db password.RegisterDB, gen password.Generator) registerRepository {
	return registerRepository{
		db: db,
	}
}

func (repo registerRepository) findLogin(user data.User) (login password.Login, err error) {
	loginSlice, err := repo.db.FilterLogin(user)
	if err != nil {
		return
	}

	if len(loginSlice) == 0 {
		err = password.ErrLoginNotFound
		return
	}

	return loginSlice[0], nil
}

func (repo registerRepository) register(user data.User, hashed password.HashedPassword) error {
	return repo.db.RegisterPassword(user, hashed)
}
