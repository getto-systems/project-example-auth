package db

import (
	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

type (
	MemoryStore struct {
		userPassword map[data.UserID]data.HashedPassword
	}
)

func (store *MemoryStore) DB() password.DB {
	return store
}

func NewPasswordStore() *MemoryStore {
	return &MemoryStore{
		userPassword: make(map[data.UserID]data.HashedPassword),
	}
}

func (store *MemoryStore) RegisterUserPassword(user data.User, password data.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

func (store *MemoryStore) FindUserPassword(user data.User) (data.HashedPassword, error) {
	password, ok := store.userPassword[user.UserID()]
	if !ok {
		return nil, errors.New("user password not found")
	}
	return password, nil
}
