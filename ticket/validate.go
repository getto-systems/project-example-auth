package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ValidateEventPublisher interface {
	ValidateTicket(data.Request)
	ValidateTicketFailed(data.Request, error)
	AuthenticatedByTicket(data.Request, data.User)
}
