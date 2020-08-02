package errors

import (
	"fmt"
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

func (err Error) TicketValidateError() bool {
	return err.action == "Ticket.Validate"
}
func (err Error) PasswordCheckError() bool {
	return err.action == "Password.Check"
}
