package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"

	"errors"
)

const (
	GENERATE_NONCE_LIMIT = 10
)

type issuer struct {
	pub        ticket.IssueEventPublisher
	signer     ticket.Signer
	expiration ticket.Expiration
	repo       issueRepository
}

func newIssuer(
	pub ticket.IssueEventPublisher,
	db ticket.IssueDB,
	signer ticket.Signer,
	expiration ticket.Expiration,
	gen ticket.NonceGenerator,
) issuer {
	return issuer{
		pub:        pub,
		signer:     signer,
		expiration: expiration,
		repo:       newIssueRepository(db, gen),
	}
}

func (issuer issuer) issue(request data.Request, user data.User) (ticket.Ticket, ticket.Nonce, data.Expires, error) {
	expires := issuer.expiration.Expires(request)
	limit := issuer.expiration.ExtendLimit(request)

	issuer.pub.IssueTicket(request, user, expires, limit)

	nonce, err := issuer.repo.register(user, expires, limit)
	if err != nil {
		issuer.pub.IssueTicketFailed(request, user, expires, limit, err)
		return nil, "", data.Expires{}, err
	}

	ticket, err := issuer.signer.Sign(nonce, user, expires)
	if err != nil {
		issuer.pub.IssueTicketFailed(request, user, expires, limit, err)
		return nil, "", data.Expires{}, err
	}

	return ticket, nonce, expires, nil
}

type issueRepository struct {
	db  ticket.IssueDB
	gen ticket.NonceGenerator
}

func newIssueRepository(db ticket.IssueDB, gen ticket.NonceGenerator) issueRepository {
	return issueRepository{
		db:  db,
		gen: gen,
	}
}

func (repo issueRepository) register(user data.User, expires data.Expires, limit data.ExtendLimit) (ticket.Nonce, error) {
	errNonceAlreadyExists := errors.New("nonce already exists")

	for count := 0; count < GENERATE_NONCE_LIMIT; count++ {
		nonce, err := repo.gen.GenerateNonce()
		if err != nil {
			return "", err
		}

		nonce, err = repo.db.RegisterTransaction(nonce, func(nonce ticket.Nonce) error {
			if repo.db.NonceExists(nonce) {
				return errNonceAlreadyExists
			}

			err := repo.db.RegisterTicket(nonce, user, expires, limit)
			if err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			return nonce, nil
		}
		if err != errNonceAlreadyExists {
			return "", err
		}
	}

	return "", errors.New("generate nonce try failed")
}
