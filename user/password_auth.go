package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type UserPasswordAuth struct {
	pub  UserPasswordAuthEventPublisher
	repo UserPasswordRepository

	request data.Request
	userID  data.UserID
}

func (user UserPasswordAuth) Authenticate(password data.RawPassword) error {
	user.pub.PasswordMatching(user.request, user.userID)

	err := user.repo.NewPassword(user.userID).Match(password)
	if err != nil {
		user.pub.PasswordMatchFailed(user.request, user.userID, err)
		return err
	}

	return nil
}

type UserPasswordAuthEventPublisher interface {
	PasswordMatching(data.Request, data.UserID)
	PasswordMatchFailed(data.Request, data.UserID, error)
}

type UserPasswordAuthEventHandler interface {
	UserPasswordAuthEventPublisher
}

type UserPasswordRepository struct {
	db      UserPasswordDB
	matcher PasswordMatcher
}

func (repo UserPasswordRepository) NewPassword(userID data.UserID) Password {
	hashed, err := repo.db.FindUserPassword(userID)
	if err != nil {
		// always fail when password not found
		return NullPassword{}
	}

	return HashedPassword{
		matcher: repo.matcher,
		hashed:  hashed,
	}
}

type UserPasswordDB interface {
	FindUserPassword(data.UserID) (data.HashedPassword, error)
}

type PasswordMatcher interface {
	MatchPassword(data.HashedPassword, data.RawPassword) error
}

type UserPasswordAuthFactory struct {
	pub  UserPasswordAuthEventPublisher
	repo UserPasswordRepository
}

func NewUserPasswordAuthFactory(pub UserPasswordAuthEventPublisher, db UserPasswordDB, matcher PasswordMatcher) UserPasswordAuthFactory {
	return UserPasswordAuthFactory{
		pub: pub,
		repo: UserPasswordRepository{
			db:      db,
			matcher: matcher,
		},
	}
}

func (f UserPasswordAuthFactory) New(request data.Request, userID data.UserID) UserPasswordAuth {
	return UserPasswordAuth{
		pub:  f.pub,
		repo: f.repo,

		request: request,
		userID:  userID,
	}
}
