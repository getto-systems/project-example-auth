package errors

import (
	"fmt"
)

type (
	Error struct {
		action   string
		message  string
		category string
	}
)

func NewError(action string, message string) error {
	return Error{
		action:  action,
		message: message,
	}
}
func NewErrorAsCategory(action string, message string, err Error) error {
	return Error{
		action:   action,
		message:  message,
		category: err.category,
	}
}
func NewCategory(category string) Error {
	return Error{
		category: category,
	}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s/%s", err.action, err.message)
}

func (err Error) Is(target error) bool {
	if len(err.category) == 0 {
		return false
	}

	t, ok := target.(Error)
	if !ok {
		return false
	}
	if len(t.category) == 0 {
		return false
	}
	return err.category == t.category
}
