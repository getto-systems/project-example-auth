package memory

import (
	"github.com/getto-systems/project-example-id/user"
)

type MemoryStore struct {
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{}
}

func (store MemoryStore) UserRoles(userID user.UserID) user.Roles {
	return nil // []string{"admin", "user"} // TODO fetch store data
}

func (store MemoryStore) MatchUserPassword(userID user.UserID, password user.Password) bool {
	return true // TODO match store data
}
