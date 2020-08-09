package log

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger interface {
		Audit(Entry)
		Error(Entry)
		Info(Entry)
		Debug(Entry)
	}

	Entry struct {
		Message string
		Request request.Request

		User  *user.User
		Login *user.Login

		Credential    *CredentialEntry
		PasswordReset *PasswordResetEntry

		Error error
	}

	CredentialEntry struct {
		ApiRoles    *credential.ApiRoles
		Expires     *expiration.Expires
		ExtendLimit *expiration.ExtendLimit
	}

	PasswordResetEntry struct {
		Session     *password_reset.Session
		Status      *password_reset.Status
		Destination *password_reset.Destination
		Expires     *expiration.Expires
	}
)
