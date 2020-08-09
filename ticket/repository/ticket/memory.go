package ticket_repository_ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/ticket/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/user"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		ticket map[credential.TicketNonce]memoryTicketData
	}

	memoryTicketData struct {
		user    user.User
		expires expiration.Expires
		limit   expiration.ExtendLimit
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ticket: make(map[credential.TicketNonce]memoryTicketData),
	}
}

func (store *MemoryStore) repo() infra.TicketRepository {
	return store
}

func (store *MemoryStore) FindUserAndExpires(nonce credential.TicketNonce) (_ user.User, _ expiration.Expires, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.user, data.expires, true, nil
}

func (store *MemoryStore) FindExtendLimit(nonce credential.TicketNonce) (_ expiration.ExtendLimit, found bool, err error) {
	data, found := store.ticket[nonce]
	if !found {
		return
	}
	return data.limit, true, nil
}

func (store *MemoryStore) UpdateExpires(nonce credential.TicketNonce, expires expiration.Expires) (err error) {
	data, found := store.ticket[nonce]
	if !found {
		return nil
	}

	data.expires = expires
	store.ticket[nonce] = data

	return nil
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
		return nil
	}

	data.expires = expiration.EmptyExpires()
	data.limit = expiration.EmptyExtendLimit()
	store.ticket[nonce] = data

	return nil
}

func (store *MemoryStore) RegisterTicket(nonceGenerator infra.TicketNonceGenerator, user user.User, expires expiration.Expires, limit expiration.ExtendLimit) (_ credential.TicketNonce, err error) {
	for count := 0; count < GENERATE_LIMIT; count++ {
		nonce, genErr := nonceGenerator.GenerateTicketNonce()
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
