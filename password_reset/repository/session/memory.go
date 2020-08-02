package session

import (
	"errors"

	"github.com/getto-systems/project-example-id/data/password_reset"
)

const (
	GENERATE_LIMIT = 10
)

type (
	MemoryStore struct {
		session map[password_reset.Session]password_reset.SessionData
		token   map[password_reset.Token]password_reset.Session
		status  map[password_reset.Session]password_reset.Status
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		session: make(map[password_reset.Session]password_reset.SessionData),
		token:   make(map[password_reset.Token]password_reset.Session),
		status:  make(map[password_reset.Session]password_reset.Status),
	}
}

func (store *MemoryStore) db() password_reset.SessionRepository {
	return store
}

func (store *MemoryStore) FindStatus(session password_reset.Session) (_ password_reset.SessionData, _ password_reset.Status, found bool, err error) {
	status, found := store.status[session]
	if !found {
		return
	}

	data, found := store.session[session]
	if !found {
		return
	}

	return data, status, found, nil
}

func (store *MemoryStore) FindSession(token password_reset.Token) (_ password_reset.SessionData, found bool, err error) {
	session, found := store.token[token]
	if !found {
		return
	}

	data, found := store.session[session]
	if !found {
		return
	}

	return data, true, nil
}

func (store *MemoryStore) RegisterSession(gen password_reset.SessionGenerator, data password_reset.SessionData) (_ password_reset.Session, _ password_reset.Token, err error) {
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

		store.session[session] = data
		store.token[token] = session
		store.status[session] = password_reset.NewStatus()

		return session, token, nil
	}

	err = errors.New("generate reset try failed")
	return
}

func (store *MemoryStore) UpdateStatusToProcessing(session password_reset.Session) (err error) {
	_, found := store.status[session]
	if !found {
		return errors.New("status not found")
	}

	// TODO ステータス変更
	//store.status[session] = status.Processing()

	return nil
}
func (store *MemoryStore) UpdateStatusToFailed(session password_reset.Session, cause error) (err error) {
	_, found := store.status[session]
	if !found {
		return errors.New("status not found")
	}

	// TODO ステータス変更
	//store.status[session] = status.Failed(cause)

	return nil
}
func (store *MemoryStore) UpdateStatusToSuccess(session password_reset.Session, dest password_reset.Destination) (err error) {
	_, found := store.status[session]
	if !found {
		return errors.New("status not found")
	}

	// TODO ステータス変更
	//store.status[session] = status.Success(dest)

	return nil
}

// for test
func (store *MemoryStore) GetResetSessionData(session password_reset.Session) (password_reset.SessionData, bool) {
	data, found := store.session[session]
	return data, found
}
