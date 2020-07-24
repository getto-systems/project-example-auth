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

func (store *MemoryStore) FindUser(login password.Login) (_ data.User, _ bool, err error) {
	userID, found := store.userID[login.ID()]
	return data.NewUser(userID), found, nil
}

func (store *MemoryStore) FindLogin(user data.User) (_ password.Login, _ bool, err error) {
	loginID, found := store.loginID[user.UserID()]
	return password.NewLogin(loginID), found, nil
}

func (store *MemoryStore) FindPassword(login password.Login) (_ data.User, _ password.HashedPassword, found bool, err error) {
	userID, found := store.userID[login.ID()]
	if !found {
		return
	}

	hashed, found := store.userPassword[userID]
	if !found {
		found = false
		return
	}

	return data.NewUser(userID), hashed, true, nil
}

func (store *MemoryStore) RegisterPassword(user data.User, password password.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

// for test
func (store *MemoryStore) GetUserPassword(user data.User) (password.HashedPassword, bool) {
	password, found := store.userPassword[user.UserID()]
	return password, found
}

// for test
func (store *MemoryStore) RegisterLogin(user data.User, login password.Login) {
	store.userID[login.ID()] = user.UserID()
	store.loginID[user.UserID()] = login.ID()
}
