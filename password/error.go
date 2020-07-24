package password

import (
	"errors"
	"fmt"
)

func newError(message string) error {
	return errors.New(fmt.Sprintf("Password/%s", message))
}
