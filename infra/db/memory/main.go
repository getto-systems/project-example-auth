package memory

import (
	"github.com/getto-systems/project-example-id/basic"
)

type MemoryStore struct {
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{}
}

func (store MemoryStore) UserRoles(userID basic.UserID) basic.Roles {
	return []string{"admin"} // TODO fetch store data
}

func (store MemoryStore) UserPassword(userID basic.UserID) (basic.HashedPassword, error) {
	return []byte("$2a$10$1HRHcllzEujLaWLcFjnXl.94GY8/Q5pu1Speo3UiPkapbq5m901ZK"), nil // TODO match store data
}
