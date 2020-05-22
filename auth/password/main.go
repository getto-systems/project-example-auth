package password

import (
	"github.com/getto-systems/project-example-id/auth"
)

type Credential struct {
	db     CredentialStore
	userID auth.UserID
}

type Password string

type CredentialStore interface {
	MatchPassword(auth.UserID, Password) bool
}

func NewCredential(db CredentialStore, userID auth.UserID) Credential {
	return Credential{db, userID}
}

func (credential Credential) Match(password Password) bool {
	return credential.db.MatchPassword(credential.userID, password)
}
