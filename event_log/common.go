package event_log

import (
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
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
	User    *data.User
	Roles   *data.Roles
	Expires *data.Expires
	Error   error
}
