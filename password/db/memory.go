package db

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"

	"errors"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		userPassword map[data.UserID]password.HashedPassword
		userID       map[password.LoginID]data.UserID
		loginID      map[data.UserID]password.LoginID

		reset       map[password.ResetID]password.ResetUser
		resetToken  map[password.ResetToken]password.ResetID
		resetStatus map[password.ResetID]password.ResetStatus
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

func (store *MemoryStore) FilterPassword(login password.Login) ([]password.Password, error) {
	userID, ok := store.userID[login.ID()]
	if !ok {
		return nil, nil
	}

	hashed, ok := store.userPassword[userID]
	if !ok {
		return nil, nil
	}

	return []password.Password{password.NewPassword(data.NewUser(userID), hashed)}, nil
}

func (store *MemoryStore) FilterLogin(user data.User) ([]password.Login, error) {
	loginID, ok := store.loginID[user.UserID()]
	if !ok {
		return nil, nil
	}

	return []password.Login{password.NewLogin(loginID)}, nil
}

func (store *MemoryStore) RegisterPassword(user data.User, password password.HashedPassword) error {
	store.userPassword[user.UserID()] = password
	return nil
}

func (store *MemoryStore) FilterUserByLogin(login password.Login) ([]data.User, error) {
	userID, ok := store.userID[login.ID()]
	if !ok {
		return nil, nil
	}

	return []data.User{data.NewUser(userID)}, nil
}

func (store *MemoryStore) RegisterReset(
	gen password.ResetGenerator,
	user data.User,
	requestedAt data.RequestedAt,
	expires data.Expires,
) (password.Reset, password.ResetToken, error) {
	for count := 0; count < GENERATE_LIMIT; count++ {
		id, token, err := gen.Generate()
		if err != nil {
			return password.Reset{}, "", err
		}

		_, ok := store.reset[id]
		if ok {
			continue
		}

		_, ok = store.resetToken[token]
		if ok {
			continue
		}

		store.reset[id] = password.NewResetUser(user, expires)
		store.resetToken[token] = id
		store.resetStatus[id] = password.NewResetStatusRequested(requestedAt.Time())

		return password.NewReset(id), token, nil
	}

	return password.Reset{}, "", errors.New("generate reset try failed")
}

func (store *MemoryStore) FilterResetStatus(reset password.Reset) ([]password.ResetStatus, error) {
	status, ok := store.resetStatus[reset.ID()]
	if !ok {
		return nil, nil
	}

	return []password.ResetStatus{status}, nil
}

func (store *MemoryStore) FilterResetUser(token password.ResetToken) ([]password.ResetUser, error) {
	id, ok := store.resetToken[token]
	if !ok {
		return nil, nil
	}

	user, ok := store.reset[id]
	if !ok {
		return nil, nil
	}

	return []password.ResetUser{user}, nil
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
