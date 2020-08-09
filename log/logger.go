package log

import (
	"github.com/getto-systems/project-example-id/misc/expiration"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type Logger interface {
	Audit(Entry)
	Error(Entry)
	Info(Entry)
	Debug(Entry)
}

type Entry struct {
	Message string
	Request request.Request

	User  *user.User
	Login *user.Login

	TicketNonce *credential.TicketNonce
	ApiRoles    *credential.ApiRoles

	CredentialExpires     *expiration.Expires
	CredentialExtendLimit *expiration.ExtendLimit

	ResetSession        *password_reset.Session
	ResetStatus         *password_reset.Status
	ResetDestination    *password_reset.Destination
	ResetSessionExpires *expiration.Expires

	Error error
}
