package password_reset_repository_destination

import (
	"github.com/getto-systems/project-example-auth/password_reset/infra"

	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	MemoryStore struct {
		destination map[user.User]password_reset.Destination
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		destination: make(map[user.User]password_reset.Destination),
	}
}

func (store *MemoryStore) db() infra.DestinationRepository {
	return store
}

func (store *MemoryStore) FindDestination(user user.User) (_ password_reset.Destination, found bool, err error) {
	dest, found := store.destination[user]
	if !found {
		return
	}

	return dest, found, nil
}

func (store *MemoryStore) RegisterDestination(user user.User, dest password_reset.Destination) (err error) {
	store.destination[user] = dest
	return nil
}
