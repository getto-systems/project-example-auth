package _usecase

import (
	"github.com/getto-systems/project-example-id/_misc/errors"
)

var (
	ErrTicketValidate = errors.NewError("Ticket.Validate", "")
	ErrPasswordCheck  = errors.NewError("Password.Check", "")
)
