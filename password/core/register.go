package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type registerer struct {
	pub  password.RegisterEventPublisher
	repo registerRepository
}

func newRegisterer(
	pub password.RegisterEventPublisher,
	db password.RegisterDB,
	gen password.Generator,
) registerer {
	return registerer{
		pub:  pub,
		repo: newRegisterRepository(db, gen),
	}
}

func (registerer registerer) register(request data.Request, user data.User, password password.RawPassword) error {
	registerer.pub.RegisterPassword(request, user)

	err := checkPassword(password)
	if err != nil {
		registerer.pub.RegisterPasswordFailed(request, user, err)
		return err
	}

	err = registerer.repo.registerPassword(user, password)
	if err != nil {
		registerer.pub.RegisterPasswordFailed(request, user, err)
		return err
	}

	registerer.pub.PasswordRegistered(request, user)

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

func (repo registerRepository) registerPassword(user data.User, password password.RawPassword) error {
	hashed, err := repo.gen.GeneratePassword(password)
	if err != nil {
		return err
	}

	return repo.db.RegisterUserPassword(user, hashed)
}
