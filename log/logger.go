package log

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password_reset"
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

	Expires     *time.Expires
	ExtendLimit *time.ExtendLimit

	ResetSession     *password_reset.Session
	ResetStatus      *password_reset.Status
	ResetDestination *password_reset.Destination

	Error error
}
