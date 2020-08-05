package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		ticket map[ticket.Nonce]memoryTicketData
	}

	memoryTicketData struct {
		user    user.User
		expires time.Expires
		limit   time.ExtendLimit
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ticket: make(map[ticket.Nonce]memoryTicketData),
	}
}

func (store *MemoryStore) repo() ticket.TicketRepository {
	return store
}

func (store *MemoryStore) FindUserAndExpires(nonce ticket.Nonce) (_ user.User, _ time.Expires, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, data.expires, true, nil
}

func (store *MemoryStore) FindUserAndExtendLimit(nonce ticket.Nonce) (_ user.User, _ time.ExtendLimit, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, data.limit, true, nil
}

func (store *MemoryStore) FindUser(nonce ticket.Nonce) (_ user.User, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, true, nil
}

func (store *MemoryStore) DeactivateExpiresAndExtendLimit(nonce ticket.Nonce) (err error) {
	data, found := store.ticket[nonce]
	if !found {
		// 見つからない場合は何もせず、特にエラーにもしない
		return nil
	}

	data.expires = time.EmptyExpires()
	data.limit = time.EmptyExtendLimit()
	store.ticket[nonce] = data

	return nil
}

func (store *MemoryStore) RegisterTicket(gen ticket.NonceGenerator, user user.User, expires time.Expires, limit time.ExtendLimit) (_ ticket.Nonce, err error) {
	for count := 0; count < GENERATE_LIMIT; count++ {
		nonce, genErr := gen.GenerateNonce()
		if genErr != nil {
			err = genErr
			return
		}

		_, found := store.ticket[nonce]
		if found {
			continue
		}

		store.ticket[nonce] = memoryTicketData{
			user:    user,
			expires: expires,
			limit:   limit,
		}

		return nonce, nil
	}

	err = errors.New("generate reset try failed")
	return
}
