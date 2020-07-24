package reset_session

import (
	"github.com/getto-systems/project-example-id/password"

	"errors"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		session map[password.ResetSessionID]password.ResetSessionData
		token   map[password.ResetToken]password.ResetSessionID
		status  map[password.ResetSessionID]password.ResetStatus
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		session: make(map[password.ResetSessionID]password.ResetSessionData),
		token:   make(map[password.ResetToken]password.ResetSessionID),
		status:  make(map[password.ResetSessionID]password.ResetStatus),
	}
}

func (store *MemoryStore) db() password.ResetSessionRepository {
	return store
}

func (store *MemoryStore) FindResetStatus(session password.ResetSession) (_ password.ResetStatus, err error) {
	status, ok := store.status[session.ID()]
	if !ok {
		err = password.ErrResetSessionNotFoundResetStatus
		return
	}

	return status, nil
}

func (store *MemoryStore) FindResetSession(token password.ResetToken) (_ password.ResetSessionData, err error) {
	id, ok := store.token[token]
	if !ok {
		err = password.ErrResetSessionNotFoundResetSession
		return
	}

	data, ok := store.session[id]
	if !ok {
		err = password.ErrResetSessionNotFoundResetSession
		return
	}

	return data, nil
}

func (store *MemoryStore) RegisterResetSession(gen password.ResetSessionGenerator, data password.ResetSessionData) (_ password.ResetSession, _ password.ResetToken, err error) {
	for count := 0; count < GENERATE_LIMIT; count++ {
		id, token, genErr := gen.GenerateSession()
		if genErr != nil {
			err = genErr
			return
		}

		_, ok := store.session[id]
		if ok {
			continue
		}

		_, ok = store.token[token]
		if ok {
			continue
		}

		store.session[id] = data
		store.token[token] = id
		store.status[id] = password.NewResetStatus()

		return password.NewResetSession(id), token, nil
	}

	err = errors.New("generate reset try failed")
	return
}

// for test
func (store *MemoryStore) GetResetSessionData(id password.ResetSessionID) (password.ResetSessionData, bool) {
	session, ok := store.session[id]
	return session, ok
}
