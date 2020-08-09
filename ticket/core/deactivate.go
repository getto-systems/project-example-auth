package ticket_core

import (
	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Deactivate(request request.Request, user user.User, ticket credential.Ticket) (err error) {
	action.logger.TryToDeactivate(request, user)

	err = action.tickets.DeactivateExpiresAndExtendLimit(ticket.Nonce())
	if err != nil {
		action.logger.FailedToDeactivate(request, user, err)
		return
	}

	action.logger.Deactivate(request, user)
	return nil
}
