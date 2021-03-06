package user_repository_user

import (
	"github.com/getto-systems/project-example-auth/user/infra"

	"github.com/getto-systems/project-example-auth/user"
)

type (
	MemoryStore struct {
		user  map[user.Login]user.User
		login map[user.User]user.Login
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		user:  make(map[user.Login]user.User),
		login: make(map[user.User]user.Login),
	}
}

func (store *MemoryStore) repo() infra.UserRepository {
	return store
}

func (store *MemoryStore) FindUser(login user.Login) (_ user.User, _ bool, err error) {
	user, found := store.user[login]
	return user, found, nil
}

func (store *MemoryStore) FindLogin(user user.User) (_ user.Login, _ bool, err error) {
	login, found := store.login[user]
	return login, found, nil
}

func (store *MemoryStore) RegisterUser(user user.User, login user.Login) error {
	store.user[login] = user
	store.login[user] = login
	return nil
}
