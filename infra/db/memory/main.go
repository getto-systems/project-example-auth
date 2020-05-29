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

func (store MemoryStore) UserPassword(userID user.UserID) user.HashedPassword {
	return []byte("$2a$10$1HRHcllzEujLaWLcFjnXl.94GY8/Q5pu1Speo3UiPkapbq5m901ZK") // TODO match store data
}
