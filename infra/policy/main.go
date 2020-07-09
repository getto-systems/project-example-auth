package policy

import (
	"github.com/getto-systems/project-example-id/data"
)

type Policy struct {
}

func NewPolicy() Policy {
	return Policy{}
}

func (policy Policy) Limit(request data.Request, roles data.Roles) data.Roles {
	return roles
}
