package credential_core

import (
	"github.com/getto-systems/project-example-auth/credential/infra"

	"github.com/getto-systems/project-example-auth/credential"
)

type (
	action struct {
		logger infra.Logger

		ticketSigner infra.TicketSigner
		ticketParser infra.TicketParser

		apiTokenSinger     infra.ApiTokenSigner
		contentTokenSigner infra.ContentTokenSigner

		apiUsers infra.ApiUserRepository
	}
)

func NewAction(
	logger infra.Logger,

	ticketSign infra.TicketSign,
	apiTokenSinger infra.ApiTokenSigner,
	contentTokenSigner infra.ContentTokenSigner,

	apiUsers infra.ApiUserRepository,
) credential.Action {
	return action{
		logger: logger,

		ticketSigner: ticketSign,
		ticketParser: ticketSign,

		apiTokenSinger:     apiTokenSinger,
		contentTokenSigner: contentTokenSigner,

		apiUsers: apiUsers,
	}
}
