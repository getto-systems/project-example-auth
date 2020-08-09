package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		ticket map[credential.TicketNonce]memoryTicketData
	}

	memoryTicketData struct {
		user         user.User
		expires      time.Expires
		expireSecond time.Second
		limit        time.ExtendLimit
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ticket: make(map[credential.TicketNonce]memoryTicketData),
	}
}

func (store *MemoryStore) repo() ticket.TicketRepository {
	return store
}

func (store *MemoryStore) FindUserAndExpires(nonce credential.TicketNonce) (_ user.User, _ time.Expires, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, data.expires, true, nil
}

func (store *MemoryStore) FindExpireSecondAndExtendLimit(nonce credential.TicketNonce) (_ time.Second, _ time.ExtendLimit, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.expireSecond, data.limit, true, nil
}

func (store *MemoryStore) FindUser(nonce credential.TicketNonce) (_ user.User, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, true, nil
}

func (store *MemoryStore) DeactivateExpiresAndExtendLimit(nonce credential.TicketNonce) (err error) {
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

func (store *MemoryStore) RegisterTicket(gen credential.TicketNonceGenerator, user user.User, expires time.Expires, expireSecond time.Second, limit time.ExtendLimit) (_ credential.TicketNonce, err error) {
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
			user:         user,
			expires:      expires,
			expireSecond: expireSecond,
			limit:        limit,
		}

		return nonce, nil
	}

	err = errors.New("generate reset try failed")
	return
}
