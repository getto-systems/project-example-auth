package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

const (
	GENERATE_NONCE_LIMIT = 10
)

type Issuer struct {
	pub        EventPublisher
	signer     Signer
	expiration Expiration
	repo       issueRepository
}

func NewIssuer(
	pub EventPublisher,
	db DB,
	signer Signer,
	expiration Expiration,
	gen NonceGenerator,
) Issuer {
	return Issuer{
		pub:        pub,
		signer:     signer,
		expiration: expiration,
		repo:       newIssueRepository(db, gen),
	}
}

func (issuer Issuer) issue(request data.Request, user data.User) (Ticket, Nonce, data.Expires, error) {
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

type NonceGenerator interface {
	GenerateNonce() (Nonce, error)
}

type issueRepository struct {
	db  DB
	gen NonceGenerator
}

func newIssueRepository(db DB, gen NonceGenerator) issueRepository {
	return issueRepository{
		db:  db,
		gen: gen,
	}
}

func (repo issueRepository) register(user data.User, expires data.Expires, limit data.ExtendLimit) (Nonce, error) {
	for count := 0; count < GENERATE_NONCE_LIMIT; count++ {
		nonce, err := repo.gen.GenerateNonce()
		if err != nil {
			return "", err
		}

		nonce, err = repo.db.RegisterTransaction(nonce, func(nonce Nonce) error {
			if repo.db.NonceExists(nonce) {
				return errors.New("nonce already exists")
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
	}

	return "", errors.New("generate nonce failed")
}
