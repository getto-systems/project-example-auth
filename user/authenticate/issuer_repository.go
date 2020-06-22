package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"
)

type IssuerRepository struct {
	db ProfileDB

	issuerFactory user.IssuerFactory
}

type ProfileDB interface {
	UserProfile(data.UserID) (data.Profile, error)
}

func (repo IssuerRepository) New(ticket data.Ticket) user.Issuer {
	return repo.issuerFactory.FromTicket(ticket)
}

func (repo IssuerRepository) Find(user user.User) user.Issuer {
	profile, err := repo.db.UserProfile(user.UserID())
	if err != nil {
		return repo.issuerFactory.New(data.Profile{})
	}

	return repo.issuerFactory.New(profile)
}

func NewIssuerRepository(db ProfileDB, issuerFactory user.IssuerFactory) IssuerRepository {
	return IssuerRepository{
		db: db,

		issuerFactory: issuerFactory,
	}
}
