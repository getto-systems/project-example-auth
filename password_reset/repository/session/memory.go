package session

import (
	"errors"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/time"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		session map[password_reset.Session]sessionData
		token   map[password_reset.Token]password_reset.Session
		closed  map[password_reset.Token]sessionData
	}

	sessionData struct {
		session password_reset.Session
		token   password_reset.Token
		data    password_reset.SessionData
		dest    password_reset.Destination
		status  password_reset.Status
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		session: make(map[password_reset.Session]sessionData),
		token:   make(map[password_reset.Token]password_reset.Session),
		closed:  make(map[password_reset.Token]sessionData),
	}
}

func (store *MemoryStore) db() password_reset_infra.SessionRepository {
	return store
}

func (store *MemoryStore) FindSessionDataAndDestinationAndStatus(session password_reset.Session) (_ password_reset.SessionData, _ password_reset.Destination, _ password_reset.Status, found bool, err error) {
	data, found := store.session[session]
	if !found {
		return
	}

	return data.data, data.dest, data.status, found, nil
}

func (store *MemoryStore) FindSession(token password_reset.Token) (_ password_reset.Session, _ password_reset.SessionData, found bool, err error) {
	session, found := store.token[token]
	if !found {
		return
	}

	data, found := store.session[session]
	if !found {
		return
	}

	return session, data.data, true, nil
}

func (store *MemoryStore) CheckClosedSessionExists(token password_reset.Token) (found bool, err error) {
	_, found = store.closed[token]
	if !found {
		return
	}

	return true, nil
}

func (store *MemoryStore) CreateSession(gen password_reset_infra.SessionGenerator, data password_reset.SessionData, dest password_reset.Destination) (_ password_reset.Session, _ password_reset.Token, err error) {
	for count := 0; count < GENERATE_LIMIT; count++ {
		id, token, genErr := gen.GenerateSession()
		if genErr != nil {
			err = genErr
			return
		}

		session := password_reset.NewSession(id)

		_, found := store.session[session]
		if found {
			continue
		}

		_, found = store.token[token]
		if found {
			continue
		}

		store.session[session] = sessionData{
			session: session,
			data:    data,
			dest:    dest,
			status:  password_reset.NewStatusWaiting(data.RequestedAt().Time()),
			token:   token,
		}
		store.token[token] = session

		return session, token, nil
	}

	err = errors.New("generate reset try failed")
	return
}

func (store *MemoryStore) CloseSession(session password_reset.Session) (err error) {
	data, found := store.session[session]
	if !found {
		err = errors.New("session not found")
		return
	}

	delete(store.session, session)
	delete(store.token, data.token)

	store.closed[data.token] = data

	return nil
}

func (store *MemoryStore) UpdateStatusToSending(session password_reset.Session, requestedAt time.RequestedAt) (err error) {
	data, found := store.session[session]
	if !found {
		err = errors.New("session not found")
		return
	}

	data.status = password_reset.NewStatusSending(requestedAt.Time())
	store.session[session] = data

	return nil
}
func (store *MemoryStore) UpdateStatusToFailed(session password_reset.Session, requestedAt time.RequestedAt, cause error) (err error) {
	data, found := store.session[session]
	if !found {
		return errors.New("session not found")
	}

	data.status = password_reset.NewStatusFailed(requestedAt.Time(), cause.Error())
	store.session[session] = data

	return nil
}
func (store *MemoryStore) UpdateStatusToComplete(session password_reset.Session, requestedAt time.RequestedAt) (err error) {
	data, found := store.session[session]
	if !found {
		return errors.New("session not found")
	}

	data.status = password_reset.NewStatusComplete(requestedAt.Time())
	store.session[session] = data

	return nil
}
