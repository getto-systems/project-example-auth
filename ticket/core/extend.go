package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type extender struct {
	pub        ticket.ExtendEventPublisher
	signer     ticket.Signer
	expiration ticket.Expiration
	repo       extendRepository
}

func newExtender(
	pub ticket.ExtendEventPublisher,
	db ticket.ExtendDB,
	signer ticket.Signer,
	expiration ticket.Expiration,
) extender {
	return extender{
		pub:        pub,
		signer:     signer,
		expiration: expiration,
		repo:       newExtendRepository(db),
	}
}

func (extender extender) extend(request data.Request, nonce ticket.Nonce, user data.User) (ticket.Ticket, data.Expires, error) {
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
	db ticket.ExtendDB
}

func newExtendRepository(db ticket.ExtendDB) extendRepository {
	return extendRepository{
		db: db,
	}
}

func (repo extendRepository) expires(nonce ticket.Nonce, user data.User, expires data.Expires) (data.Expires, error) {
	limit, err := repo.db.FindTicketExtendLimit(nonce, user)
	if err != nil {
		return data.Expires{}, err
	}

	return expires.Limit(limit), nil
}
