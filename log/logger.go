package log

import (
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
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

	Nonce    *ticket.Nonce
	ApiRoles *api_token.ApiRoles

	Expires     *time.Expires
	ExtendLimit *time.ExtendLimit

	ResetSession     *password_reset.Session
	ResetStatus      *password_reset.Status
	ResetDestination *password_reset.Destination

	Error error
}
