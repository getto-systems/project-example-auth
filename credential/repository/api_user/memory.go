package api_user

import (
	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	MemoryStore struct {
		apiRoles map[user.User]credential.ApiRoles
	}
)

func (store *MemoryStore) repo() credential.ApiUserRepository {
	return store
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		apiRoles: make(map[user.User]credential.ApiRoles),
	}
}

func (store *MemoryStore) FindApiRoles(user user.User) (_ credential.ApiRoles, found bool, err error) {
	roles, found := store.apiRoles[user]
	if !found {
		return
	}
	return roles, true, nil
}

func (store *MemoryStore) RegisterApiRoles(user user.User, roles credential.ApiRoles) (err error) {
	store.apiRoles[user] = roles
	return nil
}
