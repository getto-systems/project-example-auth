package api_user

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	MemoryStore struct {
		apiRoles map[user.User]api_token.ApiRoles
	}
)

func (store *MemoryStore) repo() api_token.ApiUserRepository {
	return store
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		apiRoles: make(map[user.User]api_token.ApiRoles),
	}
}

func (store *MemoryStore) FindApiRoles(user user.User) (_ api_token.ApiRoles, found bool, err error) {
	roles, found := store.apiRoles[user]
	if !found {
		return
	}
	return roles, true, nil
}

func (store *MemoryStore) RegisterApiRoles(user user.User, roles api_token.ApiRoles) (err error) {
	store.apiRoles[user] = roles
	return nil
}
