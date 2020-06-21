package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/basic"
)

type IssuerRepository struct {
	db ProfileDB

	issuerFactory user.IssuerFactory
}

type ProfileDB interface {
	UserProfile(basic.UserID) (basic.Profile, error)
}

func (repo IssuerRepository) New(ticket basic.Ticket) user.Issuer {
	return repo.issuerFactory.FromTicket(ticket)
}

func (repo IssuerRepository) Find(userID basic.UserID) user.Issuer {
	profile, err := repo.db.UserProfile(userID)
	if err != nil {
		return repo.issuerFactory.New(basic.Profile{})
	}

	return repo.issuerFactory.New(profile)
}

func NewIssuerRepository(db ProfileDB, issuerFactory user.IssuerFactory) IssuerRepository {
	return IssuerRepository{
		db: db,

		issuerFactory: issuerFactory,
	}
}
