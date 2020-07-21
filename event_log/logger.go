package event_log

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

type Logger interface {
	Audit(Entry)
	Info(Entry)
	Debug(Entry)
}

type Entry struct {
	Message string
	Request data.Request
	Nonce   *ticket.Nonce
	Login   *password.Login
	User    *data.User
	Roles   *data.Roles
	Expires *data.Expires
	Error   error
}
