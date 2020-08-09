package ticket

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	TicketRepository interface {
		FindUserAndExpires(credential.TicketNonce) (user.User, time.Expires, bool, error)

		FindExpireSecondAndExtendLimit(credential.TicketNonce) (time.Second, time.ExtendLimit, bool, error)
		UpdateExpires(credential.TicketNonce, time.Expires) error

		FindUser(credential.TicketNonce) (user.User, bool, error)
		DeactivateExpiresAndExtendLimit(credential.TicketNonce) error

		RegisterTicket(credential_infra.TicketNonceGenerator, user.User, time.Expires, time.Second, time.ExtendLimit) (credential.TicketNonce, error)
	}
)
