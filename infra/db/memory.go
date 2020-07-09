package db

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrUserProfileNotFound  = errors.New("user profile not found")
	ErrUserPasswordNotFound = errors.New("user password not found")
)

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (store *MemoryStore) FindUserProfile(userID data.UserID) (data.Profile, error) {
	if userID != "admin" {
		return data.Profile{}, ErrUserProfileNotFound
	}

	return data.Profile{
		UserID: userID,
		Roles:  []string{"admin"},
	}, nil
}

func (store *MemoryStore) FindUserPassword(userID data.UserID) (data.HashedPassword, error) {
	if userID != "admin" {
		return nil, ErrUserPasswordNotFound
	}

	return []byte("$2a$10$1HRHcllzEujLaWLcFjnXl.94GY8/Q5pu1Speo3UiPkapbq5m901ZK"), nil
}
