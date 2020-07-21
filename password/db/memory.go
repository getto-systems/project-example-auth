package db

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"

	"errors"
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

func (store *MemoryStore) db() password.DB {
	return store
}

func (store *MemoryStore) RegisterPasswordOfUser(user data.User, password password.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

func (store *MemoryStore) FindLoginByUser(user data.User) (password.Login, error) {
	loginID, ok := store.loginID[user.UserID()]
	if !ok {
		return password.Login{}, errors.New("login id not found")
	}

	return password.NewLogin(loginID), nil
}

func (store *MemoryStore) FindPasswordByLogin(login password.Login) (data.User, password.HashedPassword, error) {
	userID, ok := store.userID[login.ID()]
	if !ok {
		return data.User{}, nil, errors.New("user id not found")
	}

	password, ok := store.userPassword[userID]
	if !ok {
		return data.User{}, nil, errors.New("user password not found")
	}
	return data.NewUser(userID), password, nil
}

// for test
func (store *MemoryStore) GetUserPassword(user data.User) (password.HashedPassword, bool) {
	password, ok := store.userPassword[user.UserID()]
	return password, ok
}

// for test
func (store *MemoryStore) RegisterUserLogin(user data.User, login password.Login) {
	store.userID[login.ID()] = user.UserID()
	store.loginID[user.UserID()] = login.ID()
}
