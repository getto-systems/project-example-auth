package data

import (
	"fmt"
)

var (
	ErrTicketValidate = NewError("Ticket.Validate", "")
	ErrPasswordCheck  = NewError("Password.Check", "")
)

type (
	Error struct {
		action  string
		message string
	}
)

func NewError(action string, message string) error {
	return Error{
		action:  action,
		message: message,
	}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s/%s", err.action, err.message)
}

func (err Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return err.action == t.action
}
