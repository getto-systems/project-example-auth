package ticket_core

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

// user が正しいことは確認済みでなければならない
func (action action) Deactivate(request request.Request, user user.User, nonce credential.TicketNonce) (err error) {
	action.logger.TryToDeactivate(request, user)

	err = action.tickets.DeactivateExpiresAndExtendLimit(nonce)
	if err != nil {
		action.logger.FailedToDeactivate(request, user, err)
		return
	}

	action.logger.Deactivate(request, user)
	return nil
}
