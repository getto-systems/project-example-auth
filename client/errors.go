package client

import (
	"github.com/getto-systems/project-example-id/misc/errors"
)

var (
	ErrTicketValidate = errors.NewError("Ticket.Validate", "")
	ErrPasswordCheck  = errors.NewError("Password.Check", "")
)
