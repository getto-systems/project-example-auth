package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type Extender struct {
	pub        EventPublisher
	signer     Signer
	expiration Expiration
	repo       extendRepository
}

func NewExtender(
	pub EventPublisher,
	db DB,
	signer Signer,
	expiration Expiration,
) Extender {
	return Extender{
		pub:        pub,
		signer:     signer,
		expiration: expiration,
		repo:       newExtendRepository(db),
	}
}

func (extender Extender) extend(request data.Request, nonce Nonce, user data.User) (Ticket, data.Expires, error) {
	expires := extender.expiration.Expires(request)
	expires, err := extender.repo.expires(nonce, user, expires)

	extender.pub.ExtendTicket(request, nonce, user, expires)

	if err != nil {
		extender.pub.ExtendTicketFailed(request, nonce, user, expires, err)
		return nil, data.Expires{}, err
	}

	ticket, err := extender.signer.Sign(nonce, user, expires)
	if err != nil {
		extender.pub.ExtendTicketFailed(request, nonce, user, expires, err)
		return nil, data.Expires{}, err
	}

	return ticket, expires, nil
}

type extendRepository struct {
	db DB
}

func newExtendRepository(db DB) extendRepository {
	return extendRepository{
		db: db,
	}
}

func (repo extendRepository) expires(nonce Nonce, user data.User, expires data.Expires) (data.Expires, error) {
	limit, err := repo.db.FindTicketExtendLimit(nonce, user)
	if err != nil {
		return data.Expires{}, err
	}

	return expires.Limit(limit), nil
}
