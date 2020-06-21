package policy

import (
	"github.com/getto-systems/project-example-id/policy"

	"github.com/getto-systems/project-example-id/data"
)

type PolicyChecker struct {
}

func NewPolicyChecker() PolicyChecker {
	return PolicyChecker{}
}

func (PolicyChecker) HasEnoughPermission(ticket data.Ticket, request data.Request, resource data.Resource) error {
	path := policy.Path(resource.Path)
	roles := policy.Roles(ticket.Profile.Roles)
	return policy.NewPolicy(path).Correct(roles)
}
