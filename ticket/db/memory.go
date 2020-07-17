package db

import (
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

type (
	MemoryStore struct {
		ticket    map[string]memoryTicketData
		userRoles map[data.UserID]data.Roles
	}

	memoryTicketData struct {
		userID  data.UserID
		expires data.Expires
		limit   data.ExtendLimit
	}
)

func (store *MemoryStore) DB() ticket.DB {
	return store
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ticket:    make(map[string]memoryTicketData),
		userRoles: make(map[data.UserID]data.Roles),
	}
}

func (store *MemoryStore) FindUserRoles(user data.User) (data.Roles, error) {
	roles, ok := store.userRoles[user.UserID()]
	if !ok {
		return nil, errors.New("user role not found")
	}
	return roles, nil
}

func (store *MemoryStore) FindTicketExtendLimit(nonce ticket.Nonce, user data.User) (data.ExtendLimit, error) {
	ticketData, ok := store.ticket[string(nonce)]
	if !ok {
		return data.ExtendLimit{}, errors.New("ticket extend limit not found")
	}

	if ticketData.userID != user.UserID() {
		return data.ExtendLimit{}, errors.New("ticket extend limit invalid: different user")
	}

	return ticketData.limit, nil
}

func (store *MemoryStore) RegisterTransaction(nonce ticket.Nonce, cb func(ticket.Nonce) error) (ticket.Nonce, error) {
	err := cb(nonce)
	if err != nil {
		return "", err
	}
	return nonce, nil
}

func (store *MemoryStore) RegisterTicket(nonce ticket.Nonce, user data.User, expires data.Expires, limit data.ExtendLimit) error {
	store.ticket[string(nonce)] = memoryTicketData{
		userID:  user.UserID(),
		expires: expires,
		limit:   limit,
	}
	return nil
}

func (store *MemoryStore) NonceExists(nonce ticket.Nonce) bool {
	_, ok := store.ticket[string(nonce)]
	return ok
}

func (store *MemoryStore) TicketExists(nonce ticket.Nonce, user data.User) bool {
	ticketData, ok := store.ticket[string(nonce)]
	if !ok {
		return false
	}
	return ticketData.userID == user.UserID()
}

func (store *MemoryStore) ShrinkTicket(nonce ticket.Nonce) error {
	ticketData, ok := store.ticket[string(nonce)]
	if !ok {
		return errors.New("shrink ticket failed: nonce not exists")
	}

	ticketData.limit = data.ExtendLimit{}
	store.ticket[string(nonce)] = ticketData

	return nil
}
