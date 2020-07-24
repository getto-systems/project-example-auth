package password

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type (
	MemoryStore struct {
		userPassword map[data.UserID]password.HashedPassword
		userID       map[password.LoginID]data.UserID
		loginID      map[data.UserID]password.LoginID
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		userPassword: make(map[data.UserID]password.HashedPassword),
		userID:       make(map[password.LoginID]data.UserID),
		loginID:      make(map[data.UserID]password.LoginID),
	}
}

func (store *MemoryStore) db() password.PasswordRepository {
	return store
}

func (store *MemoryStore) FindUser(login password.Login) (_ data.User, err error) {
	userID, ok := store.userID[login.ID()]
	if !ok {
		err = password.ErrPasswordNotFoundUser
		return
	}

	return data.NewUser(userID), nil
}

func (store *MemoryStore) FindLogin(user data.User) (_ password.Login, err error) {
	loginID, ok := store.loginID[user.UserID()]
	if !ok {
		err = password.ErrPasswordNotFoundLogin
		return
	}

	return password.NewLogin(loginID), nil
}

func (store *MemoryStore) FindPassword(login password.Login) (_ data.User, _ password.HashedPassword, err error) {
	userID, ok := store.userID[login.ID()]
	if !ok {
		err = password.ErrPasswordNotFoundPassword
		return
	}

	hashed, ok := store.userPassword[userID]
	if !ok {
		err = password.ErrPasswordNotFoundPassword
		return
	}

	return data.NewUser(userID), hashed, nil
}

func (store *MemoryStore) RegisterPassword(user data.User, password password.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

// for test
func (store *MemoryStore) GetUserPassword(user data.User) (password.HashedPassword, bool) {
	password, ok := store.userPassword[user.UserID()]
	return password, ok
}

// for test
func (store *MemoryStore) RegisterLogin(user data.User, login password.Login) {
	store.userID[login.ID()] = user.UserID()
	store.loginID[user.UserID()] = login.ID()
}
