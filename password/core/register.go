package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type registerer struct {
	logger password.RegisterLogger
	repo   registerRepository
}

func newRegisterer(
	logger password.RegisterLogger,
	db password.RegisterDB,
	gen password.Generator,
) registerer {
	return registerer{
		logger: logger,
		repo:   newRegisterRepository(db, gen),
	}
}

func (registerer registerer) getLogin(request data.Request, user data.User) (password.Login, error) {
	registerer.logger.TryToGetLogin(request, user)

	login, err := registerer.repo.findLogin(user)
	if err != nil {
		registerer.logger.FailedToGetLogin(request, user, err)
		return password.Login{}, err
	}

	return login, nil
}

func (registerer registerer) register(request data.Request, user data.User, password password.RawPassword) error {
	registerer.logger.TryToRegister(request, user)

	err := registerer.repo.register(user, password)
	if err != nil {
		registerer.logger.FailedToRegister(request, user, err)
		return err
	}

	registerer.logger.Registered(request, user)

	return nil
}

type registerRepository struct {
	db  password.RegisterDB
	gen password.Generator
}

func newRegisterRepository(db password.RegisterDB, gen password.Generator) registerRepository {
	return registerRepository{
		db:  db,
		gen: gen,
	}
}

func (repo registerRepository) findLogin(user data.User) (password.Login, error) {
	loginSlice, err := repo.db.FilterLogin(user)
	if err != nil {
		return password.Login{}, err
	}

	if len(loginSlice) == 0 {
		return password.Login{}, password.ErrLoginNotFound
	}

	return loginSlice[0], nil
}

func (repo registerRepository) register(user data.User, password password.RawPassword) error {
	err := checkPassword(password)
	if err != nil {
		return err
	}

	hashed, err := repo.gen.GeneratePassword(password)
	if err != nil {
		return err
	}

	return repo.db.RegisterPassword(user, hashed)
}
