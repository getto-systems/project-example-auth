package password_repository_password

import (
	"github.com/getto-systems/project-example-auth/password/infra"

	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/user"
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

func (store *MemoryStore) repo() infra.PasswordRepository {
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
