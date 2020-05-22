package memory

import (
	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/auth/password"
)

type MemoryStore struct {
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{}
}

func (store MemoryStore) UserRoles(userID auth.UserID) auth.Roles {
	return nil // TODO fetch store data
}

func (store MemoryStore) MatchPassword(userID auth.UserID, password password.Password) bool {
	return true // TODO match store data
}
