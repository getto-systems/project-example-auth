package db

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"

	"errors"
)

type (
	MemoryStore struct {
		userPassword map[data.UserID]password.HashedPassword
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		userPassword: make(map[data.UserID]password.HashedPassword),
	}
}

func (store *MemoryStore) db() password.DB {
	return store
}

func (store *MemoryStore) RegisterUserPassword(user data.User, password password.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

func (store *MemoryStore) FindUserPassword(user data.User) (password.HashedPassword, error) {
	password, ok := store.userPassword[user.UserID()]
	if !ok {
		return nil, errors.New("user password not found")
	}
	return password, nil
}

// for test
func (store *MemoryStore) GetUserPassword(user data.User) (password.HashedPassword, bool) {
	password, ok := store.userPassword[user.UserID()]
	return password, ok
}
