package credential_core

import (
	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/credential"
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

	Expires struct {
		ticket expiration.ExpireSecond
		token  expiration.ExpireSecond
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
