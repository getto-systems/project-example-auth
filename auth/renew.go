package auth

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/basic"

	"fmt"
)

type RenewAuthenticator interface {
	Authenticator
	UserFactory() user.UserFactory
}

type RenewParam struct {
	RequestedAt basic.RequestedAt

	Ticket basic.Ticket
	Path   basic.Path
}

func (param RenewParam) String() string {
	return fmt.Sprintf(
		"RenewParam{RequestedAt: %s, Ticket:%s, Path:%s}",
		param.RequestedAt.String(),
		param.Ticket,
		param.Path,
	)
}

func Renew(authenticator RenewAuthenticator, param RenewParam) (basic.Ticket, error) {
	logger := authenticator.Logger()

	logger.Debugf("renew token: %v", param)

	ticket := param.Ticket

	info, err := user.NewTicketInfo(ticket, param.Path)
	if err != nil {
		logger.Debugf("ticket check failed: %s; %v", err, param)
		return basic.Ticket{}, err
	}

	if !info.IsRenewRequired(param.RequestedAt) {
		return ticket, nil
	}

	logger.Debugf("renew token: %v/%s", param)

	user := authenticator.UserFactory().New(ticket.UserID)

	new_ticket, err := user.NewTicket(param.Path, param.RequestedAt)
	if err != nil {
		logger.Auditf("access denied: %s; %v; %v", err, param)
		return basic.Ticket{}, ErrUserAccessDenied
	}

	logger.Auditf("token renewed: %v; %s", new_ticket, param.Path)

	return new_ticket, nil
}
