package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

const (
	GENERATE_NONCE_LIMIT = 10
)

type Issuer struct {
	pub        issueEventPublisher
	signer     Signer
	expiration Expiration
	repo       issueRepository
}

type issueEventPublisher interface {
	IssueTicket(data.Request, data.User, data.Expires, data.ExtendLimit)
	IssueTicketFailed(data.Request, data.User, data.Expires, data.ExtendLimit, error)
}

type issueDB interface {
	RegisterTransaction(Nonce, func(Nonce) error) (Nonce, error)
	RegisterTicket(Nonce, data.User, data.Expires, data.ExtendLimit) error
	NonceExists(Nonce) bool
}

type NonceGenerator interface {
	GenerateNonce() (Nonce, error)
}

func NewIssuer(
	pub issueEventPublisher,
	db issueDB,
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

type issueRepository struct {
	db  issueDB
	gen NonceGenerator
}

func newIssueRepository(db issueDB, gen NonceGenerator) issueRepository {
	return issueRepository{
		db:  db,
		gen: gen,
	}
}

func (repo issueRepository) register(user data.User, expires data.Expires, limit data.ExtendLimit) (Nonce, error) {
	errNonceAlreadyExists := errors.New("nonce already exists")

	for count := 0; count < GENERATE_NONCE_LIMIT; count++ {
		nonce, err := repo.gen.GenerateNonce()
		if err != nil {
			return "", err
		}

		nonce, err = repo.db.RegisterTransaction(nonce, func(nonce Nonce) error {
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
