package policy

import (
	"fmt"
	"strings"

	"errors"
)

var (
	ErrTicketHasNotEnoughPermission = errors.New("ticket has not enough permission")
)

const (
	super_role = "admin"
)

type Path string
type Roles []string

type Policy struct {
	path Path
}

func NewPolicy(path Path) Policy {
	return Policy{
		path: path,
	}
}

func (policy Policy) Correct(roles Roles) error {
	for _, role := range roles {
		if role == super_role {
			return nil
		}

		if strings.HasPrefix(fmt.Sprintf("/%s/", role), string(policy.path)) {
			return nil
		}
	}

	return ErrTicketHasNotEnoughPermission
}
