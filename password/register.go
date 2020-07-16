package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type Register struct {
	pub  registerEventPublisher
	repo registerRepository
}

type registerEventPublisher interface {
	ValidatePassword(data.Request, data.User)
	ValidatePasswordFailed(data.Request, data.User, error)
	PasswordRegistered(data.Request, data.User)
}

func NewRegister(
	pub registerEventPublisher,
	db registerDB,
	gen Generator,
) Register {
	return Register{
		pub:  pub,
		repo: newRegisterRepository(db, gen),
	}
}

func (register Register) register(request data.Request, user data.User, password data.RawPassword) error {
	register.pub.ValidatePassword(request, user)

	err := register.repo.registerPassword(user, password)
	if err != nil {
		register.pub.ValidatePasswordFailed(request, user, err)
		return err
	}

	register.pub.PasswordRegistered(request, user)

	return nil
}

type registerRepository struct {
	db  registerDB
	gen Generator
}

type registerDB interface {
	RegisterUserPassword(data.User, data.HashedPassword) error
}

func newRegisterRepository(db registerDB, gen Generator) registerRepository {
	return registerRepository{
		db:  db,
		gen: gen,
	}
}

type Generator interface {
	GeneratePassword(data.RawPassword) (data.HashedPassword, error)
}

func (repo registerRepository) registerPassword(user data.User, password data.RawPassword) error {
	hashed, err := repo.gen.GeneratePassword(password)
	if err != nil {
		return err
	}

	return repo.db.RegisterUserPassword(user, hashed)
}
