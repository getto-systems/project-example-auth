package db

import (
	"github.com/getto-systems/project-example-id/data"
)

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (store *MemoryStore) UserProfile(userID data.UserID) (data.Profile, error) {
	return data.Profile{
		UserID: userID,
		Roles:  []string{"admin"}, // TODO fetch store data
	}, nil
}

func (store *MemoryStore) UserPassword(userID data.UserID) (data.HashedPassword, error) {
	return []byte("$2a$10$1HRHcllzEujLaWLcFjnXl.94GY8/Q5pu1Speo3UiPkapbq5m901ZK"), nil // TODO match store data
}
