package credential

import (
	infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type IssueTicket struct {
	logger infra.IssueTicketLogger
	signer infra.TicketSigner
}

func NewIssueTicket(logger infra.IssueTicketLogger, signer infra.TicketSigner) IssueTicket {
	return IssueTicket{
		logger: logger,
		signer: signer,
	}
}

func (action IssueTicket) Issue(request request.Request, user user.User, nonce credential.TicketNonce, expires time.Expires) (_ credential.Ticket, err error) {
	action.logger.TryToIssueTicket(request, user, nonce, expires)

	signature, err := action.signer.Sign(user, nonce, expires)
	if err != nil {
		action.logger.FailedToIssueTicket(request, user, nonce, expires, err)
		return
	}

	action.logger.IssueTicket(request, user, nonce, expires)
	return credential.NewTicket(signature, nonce), nil
}