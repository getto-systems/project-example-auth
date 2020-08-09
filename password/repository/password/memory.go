package password

import (
	password_infra "github.com/getto-systems/project-example-id/infra/password"

	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	MemoryStore struct {
		password map[user.User]password.HashedPassword
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		password: make(map[user.User]password.HashedPassword),
	}
}

func (store *MemoryStore) repo() password_infra.PasswordRepository {
	return store
}

func (store *MemoryStore) FindPassword(user user.User) (_ password.HashedPassword, found bool, err error) {
	hashed, found := store.password[user]
	if !found {
		found = false
		return
	}

	return hashed, true, nil
}

func (store *MemoryStore) ChangePassword(user user.User, hashed password.HashedPassword) error {
	store.password[user] = hashed
	return nil
}
