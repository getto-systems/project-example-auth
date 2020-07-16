package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

type Verifier struct {
	pub    EventPublisher
	signer Signer
}

func NewVerifier(pub EventPublisher, signer Signer) Verifier {
	return Verifier{
		pub:    pub,
		signer: signer,
	}
}

func (verifier Verifier) verify(request data.Request, ticket Ticket, nonce Nonce) (data.User, error) {
	verifier.pub.VerifyTicket(request)

	verifiedNonce, user, expires, err := verifier.signer.Verify(ticket)
	if err != nil {
		verifier.pub.VerifyTicketFailed(request, err)
		return data.User{}, err
	}

	if verifiedNonce != nonce {
		err = errors.New("ticket nonce different")
		verifier.pub.VerifyTicketFailed(request, err)
		return data.User{}, err
	}

	if request.RequestedAt().Expired(expires) {
		err = errors.New("ticket already expired")
		verifier.pub.VerifyTicketFailed(request, err)
		return data.User{}, err
	}

	verifier.pub.AuthenticatedByTicket(request, user)

	return user, nil
}
